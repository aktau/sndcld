package snd

// {
// 	"id": 2511,
// 	"kind": "user",
// 	"permalink": "moullinex",
// 	"username": "Moullinex",
// 	"uri": "http://api.soundcloud.com/users/2511",
// 	"permalink_url": "http://soundcloud.com/moullinex",
// 	"avatar_url": "http://i1.sndcdn.com/avatars-000026871864-c10oaq-large.jpg?2aaad5e"
// }
type User struct {
	Id           int64  `json:"id"`
	Kind         string `json:"kind"` // have only encountered "user" so far
	Username     string `json:"username"`
	Uri          string `json:"uri"`
	PermalinkUrl string `json:"permalink_url"`
	AvatarUrl    string `json:"avatar_url"`
}
