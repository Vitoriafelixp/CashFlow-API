CREATE TABLE entries (
	id 			INT AUTO_INCREMENT PRIMARY KEY,
    description VARCHAR(200) NOT NULL,
    amount		DECIMAL(10,2) NOT NULL CHECK (amount > 0),
    type		VARCHAR (25) NOT NULL CHECK (type IN ('receivable', 'payable')),
    due_date	DATE,
    expect_receipt_date DATE,
    category 	VARCHAR(50),
    status 		VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('paid', 'pending')),
    paid_at 	DATE,
    confirmed	BOOLEAN,
    created_at	TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CHECK (
		(type = 'payable' AND due_date IS NOT NULL AND expect_receipt_date IS NULL)
		OR 
		(type = 'receivable' AND expect_receipt_date IS NOT NULL AND due_date IS NULL)
        )
	);
    
    CREATE TABLE settings (
    id 	INT AUTO_INCREMENT PRIMARY KEY,
    restock_budget	DECIMAL(10,2) NOT NULL DEFAULT 0
    );