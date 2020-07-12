package poohchanbot

func replyRandomMessage(string) *sendingMessageQueue {
	messages := [][]string{
		{"ぷぅちゃん！"},
		{"ぷぅちゃん♪"},
		{"ぷぅちゃん♡"},
		{"♡"},
		{"♪"},
		{"ぷぅ？"},
		{"ぷっぷぷ〜ぷっぷぷ〜ぷっぷぷっぷぷぅ〜♪"},
		{"シャーッ！！！"},
		{"ギャーギャー！ギャーギャー！"},
	}

	return generateRandomTextMessageQueue(messages)
}
