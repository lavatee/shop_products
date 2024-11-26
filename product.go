package products

type Product struct {
	Name        string `db:"name"`
	Amount      int    `db:"amount"`
	Price       int    `db:"price"`
	Description string `db:"description"`
	Category    string `db:"category"`
	UserId      int    `db:"user_id"`
	Id          int    `db:"id"`
}
