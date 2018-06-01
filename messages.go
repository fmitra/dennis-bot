package main

var MessageMap = map[string][]string{
	"default": []string{
		"Hi! Tell Dennis what you want to do!",

		"What are you tracking? You can say something like " +
			"2000JPY for cornerstore sushi. Not in Japan? No problemmm " +
			"you can use any currency!",

		"Let's get started! You can say something like " +
			"4USD for coffee yesterday",

		"Dennis Dennis Dennis Dennis",

		"Hiiiiiiiii I'm Dennis!",

		"How much did you spend? You can say something like 20USD " +
			"for Dinner",
	},
	"tracking_success": []string{
		"Ok got it!",

		"Roger that!",
	},
	"tracking_error": []string{
		"I didn't understand that. You need to tell me exactly what " +
			"what your expense is. For example '0.00012BTC for Rent'",

		"I didn't get that. Try saying 'How much did I spend this week?'",
	},
	"period_total_success": []string{
		"You spent {{var}}",
	},
}
