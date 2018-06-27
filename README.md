# qvo-go-client

## Description

Unofficial Go client for the QVO payment service. All objects and calls from their [REST API](https://docs.qvo.cl/) are implemented.

## Tests

Right now only customer tests are available as that's enough for my use case, but others will be added soon.

## Documentation

The package implements everything as described on QVO's docs. Please refer to thir offical [REST API](https://docs.qvo.cl/) documentation, but also check the package's [godocs](https://godoc.org/github.com/iegomez/qvo-go-client) for more details.

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

## License

The QVO Go client is distributed under the MIT license. See also LICENSE.