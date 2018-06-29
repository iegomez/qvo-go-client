# qvo-go-client

## Description

Unofficial Go client for the [QVO payment service](https://www.qvo.cl/).  
All objects and calls from their [REST API](https://docs.qvo.cl/) are implemented.

## Documentation

The package implements everything as described on QVO's docs.   
Please refer to their offical [REST API](https://docs.qvo.cl/) documentation for information about expected parameters.  
Also check the package's [godocs](https://godoc.org/github.com/iegomez/qvo-go-client) for details about the implementation.

## Requirements

This project depends on 3 Go packages:

github.com/pkg/errors for better error handling.  
github.com/smartystreets/goconvey/convey for testing.  
github.com/sirupsen/logrus for logging.  

You may install them easilly by running this:

```
make requirements
```

Get the package with:

```
go get github.com/iegomez/qvo-go-client
```

## Tests

In order to test the package, you need to set the QVO_TEST_TOKEN env var with your sandbox api token. Just export the var in the terminal before running the tests, or add the export to your .profile, .bash_profile, .bash_rc, etc., depending on your system, and then source the file before running tests.

Only customer and plan tests (and events listing in customer test) are available, as transaction, subscription, payment, withdrawal and webpay process need real card data or user actions to be tested.   
So please file an issue for any bug you may encounter using them and I'll fix it as soon as possible.

You may run the tests like this:

```
make test
```

If you want tests to stop at first fail, run them like this:

```
mate test-fast
```

They run with debug log level. Just delete the line setting the level at the test files to run with info level, or set your level of preference.

## Usage 

After importing it, the package qvo is exposed:

```go
import "github.com/iegomez/qvo-go-client"
```

The Client expects a JWT authorization token for the API, and a sandbox/production mode bool (true for sandbox). So, using your `token`, you may intialize a pointer to a Client and then call any exported function passing the pointer:

```go
c := qvo.NewClient("your-api-token", true) //NewClient returns a pointer to a qvo client.

var where = make(map[string]map[string]interface{}) //Create map for the filters
where["name"] = make(map[string]interface{})
where["name"]["like"] = "%Test%"

plans, err := qvo.ListPlans(c, 0, 0, where, "") //To omit pages, perPage or order parameters, just pass Go's zero values for ints and string.

```

Client's default log level is Info, but youy may change it with the method SetLogLevel:

```go
c := qvo.NewClient("your-api-token", true) //NewClient returns a pointer to a qvo client.
c.SetLogLevel(log.DebugLevel)
```

## Example

Here´s a stripped example used in a real project showing a function to start a webpay transaction and another to check the transaction's status:

```go

import "github.com/iegomez/qvo-go-client"

//StartPayment receives a contract id and calls to QVO for a webpay request. It returns the qvo redirect url or an error.
//It looks for a customer at QVO's end with the client's email. If there's none, it creates it and then makes the request.
func StartPayment(db *sqlx.DB, qc *qvo.Client, rp *redis.Pool, contractID int64) (string, error) {

	//Get redis connection.
	conn := rp.Get()
	defer conn.Close()

	contract, err := GetContract(db, contractID)
	if err != nil {
		return "", err
	}

	user, err := GetInternalUser(db, contract.ClientID)
	if err != nil {
		return "", err
	}

	client, err := GetUserClient(db, user.ID)
	if err != nil {
		return "", err
	}

	var customer qvo.Customer
	log.Debugf("Paying contract with user %v", user)
	//Check that the user has a qvo customer id. If not, try to create a customer.
	if user.QvoID == "" {
		log.Debugf("trying to create customer for user %d", user.ID)
		customer, err = qvo.CreateCustomer(qc, fmt.Sprintf("%s %s", client.Name, client.Lastname), user.Email)
		if err != nil {
			return "", err
		}
		//Set the customer id for the user.
		err = SetQvoID(db, user.ID, customer.ID)
		if err != nil {
			return "", err
		}
	} else {
		log.Debugf("trying to retrieve customer for user %d with qvo_id %s", user.ID, user.QvoID)
		customer, err = qvo.GetCustomer(qc, user.QvoID)
		if err != nil {
			return "", err
		}
	}

	//Make the transaction request.
	resp, err := qvo.WebpayTransaction(qc, customer.ID, fmt.Sprintf("%s#/webpay_return", common.Host), fmt.Sprintf("user %d (%s) attempts to pay contract %d", user.ID, user.Email, contract.ID), contract.Price)
	if err != nil {
		return "", err
	}

	//Now we should store the transaction id in redis so it may expire. It should contain the contract id to mark it as paid if everything goes ok.
	reply, err := conn.Do("SETEX", resp.TransactionID, int(time.Until(resp.ExpirationDate).Seconds()), contract.ID)
	if err != nil {
		log.Errorf("couldn't set transaction for contract id %d: %s", contract.ID, err)
		return "", err
	}

	log.Debugf("transaction was set with reply: %v", reply)

	return resp.RedirectURL, nil
}



//CheckPayment checks the status of a given transaction. If it's ok, it sets the contract id as paid. If not, it deletes it.
func CheckPayment(db *sqlx.DB, qc *qvo.Client, rp *redis.Pool, transactionID string) (bool, error) {
	//Get the contract id from redis.
	conn := rp.Get()
	defer conn.Close()

	contractID, redisErr := redis.String(conn.Do("GET", transactionID))
	if redisErr != nil {
		log.Errorf("couldn't get contract for transaction %s: %s", transactionID, redisErr)
		return false, redisErr
	}

	cID, err := strconv.ParseInt(contractID, 10, 64)
	if err != nil {
		log.Errorf("strconv error: %s", err)
		return false, err
	}

	//Now check against QVO that the transaction is ok.
	transaction, err := qvo.GetTransaction(qc, transactionID)
	if err != nil {
		log.Errorf("get transaction error: %s", err)
		return false, err
	}

	//We have a transaction, let's check the status.
	//If successful
	if transaction.Status == qvo.Successful {

		payErr := PayContract(db, cID, transaction)
		if payErr != nil {
			log.Errorf("pay contract error: %s", err)
			return false, payErr
		}
		//Everything's ok, delete the transaction ID from redis.
		resp, err := conn.Do("DEL", transactionID)
		if err != nil {
			log.Errorf("Couldn't delete key %d, err: %s", transactionID, err)
		}
		log.Debugf("redis delete key: %s", resp)
		return true, nil
	} else if transaction.Status == qvo.Waiting {
		//If transaction isn't ready, the frontend should keep asking, so signal that.
		log.Errorf("sending retry: %s", err)
		return false, errors.New("retry")
	}
	//On any other status, return false but with a new error to signal the client why it failed.
	//Also, delete the redis key and the contract.
	resp, err := conn.Do("DEL", transactionID)
	if err != nil {
		log.Errorf("Couldn't delete key %d, err: %s", transactionID, err)
	}
	log.Debugf("redis delete key: %s", resp)

	err = DeleteContract(db, cID)
	if err != nil {
		log.Errorf("Couldn't delete contract %d, err: %s", contractID, err)
	}

	return false, errors.Errorf("transaction error: %s", transaction.Status)
}

```

## Caveats

You should be careful about some caveats with the original API:

For list filters and order strings you should really check the offical docs for their syntax, as qvo errors won't mention any issue in the param field. For example, if you pass a random field name as filter, or you mistype the order (e.g., `created ASC` instead of the correct `created_at ASC`), QVO will respond with a 500 status code.

The API is somewhat inconsistent with int an decimal fields. First, it allows to pass an int or a string which contains an int as the `price` field of a plan with CLP currency, but won't allow a float nor a string containing a float. Oddly enough, on creation or retrieval, it'll return a float string for the same field. So you may create a plan with price 19000 or "19000" if the currency is CLP (UF allows both ints and floats), but not 19000.0 or "19000.0", and the API will return it with "19000.0" (always a string, never 19000.0) as the `price`. I could deal with this at the client implementation, but it seems messy and I've already reported it, so hopefully it'll be addressed soon.

I'll update this section if there's any change on the API.

## Contributing

Report any bug by filing an issue.

There's not much to be added client wise, as it is faithful to the public docs from QVO. Nevertheless, please file issues with the `enhancement` or `feature` tag for anything that's not included and you'd like to see implemented (stats, comparisons between customer, or any other thing you can think of).

Of course, feel free to submit a PR for any of the above.

## License

The QVO Go client is distributed under the MIT license. See also LICENSE.