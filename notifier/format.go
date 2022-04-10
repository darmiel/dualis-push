package notifier

var DefaultFormat = Format{
	NewGradeMessageBody:  "🎉 In [%course%] hast du %grade% Pommes erhalten.",
	NewGradeMessageTitle: "🫣 New Grade arrived!",
}

type Format struct {
	NewGradeMessageTitle string
	NewGradeMessageBody  string
}
