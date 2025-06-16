package products

type Product struct {
	Name        string `db:"name" json:"name"`
	Amount      int    `db:"amount" json:"amount"`
	Price       int    `db:"price" json:"price"`
	Description string `db:"description" json:"description"`
	Category    string `db:"category" json:"category"`
	UserId      int    `db:"user_id" json:"userId"`
	Id          int    `db:"id" json:"id"`
	PhotoUrl    string `db:"photo_url" json:"photoUrl"`
}
