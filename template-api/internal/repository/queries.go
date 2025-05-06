package repository

const (
	createProductQuery = `
	INSERT INTO
	    products (
			name,
		    description,
		    price,
		    stock
		) VALUES ($1, $2, $3, $4) RETURNING id
`
	getProductByIdQuery = `
	SELECT
    	id,
    	name,
    	description,
    	price,
    	stock
	FROM products WHERE id = $1
`
	getAllProductsQuery = `
	SELECT 
	    id,
	    name,
	    description,
	    price,
	    stock
	FROM products
`
)
