package main

// Creates an expense item from a Wit.ai response
func createExpense(witResponse WitResponse) {
	date := witResponse.getDate()
	amount, currency, _ := witResponse.getAmount()
	description, _ := witResponse.getDescription()

	Db.Create(&Expense{
		Date: date,
		Description: description,
		Total: amount,
		Currency: currency,
	})
}
