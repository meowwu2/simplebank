package db

import (
	"context"
	"database/sql"
)

// Store provides all functions to excute db queries and transactions
type Store interface{
	Querier
	TransferTx(ctx context.Context,arg TransferTxParams)(TransferTxResult,error)

}
// Store provides all functions to excute SQL queries and transactions
type SQLStore struct{
	*Queries
	db *sql.DB
}

// NewStore create a new Store
func NewStore(db *sql.DB) Store{
 return &SQLStore{
	Queries: New(db),
	db: db,
 }
}


// execTx excutes a function within a database transaction  
func (store *SQLStore) execTx(ctx context.Context,fn func(*Queries) error) error{
	tx,err := store.db.BeginTx(ctx,nil)
	if err != nil{
		return err
	}
	q:=New(tx)
	err = fn(q)
	if err!=nil{
		if rbErr := tx.Rollback();rbErr!=nil{
		}
		return err
	}
	return tx.Commit()

}

type TransferTxParams struct{
	FromAccountID int64     `json:"from_account_id"`
	ToAccountID   int64     `json:"to_account_id"`
	Amount        int64     `json:"amount"`
}

// the result of transfer transaction
type TransferTxResult struct{
	Transfer	Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntries Entry    `json:"from_entry"`
	ToEntries   Entry    `json:"to_entry"`
}
// TransferTx performs a money transfer from one to the other
// It creates a trasfer record, add account entries and update accounts' balance within a single database transaction
func (store *SQLStore)TransferTx(ctx context.Context,arg TransferTxParams)(TransferTxResult,error){
	var result TransferTxResult
	err :=store.execTx(ctx,func(q *Queries)error{
		var err error

		
		result.Transfer,err = q.CreateTransfer(ctx,CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err!=nil{
			return err
		}

		result.FromEntries,err = q.CreateEntry(ctx,CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err!=nil{
			return err
		}

		result.ToEntries,err = q.CreateEntry(ctx,CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err!=nil{
			return err
		}
		if arg.FromAccountID<arg.ToAccountID{
			result.FromAccount,result.ToAccount,err=addMony(ctx,q,arg.FromAccountID,-arg.Amount,arg.ToAccountID,arg.Amount)
		}else{
			result.ToAccount,result.FromAccount,err=addMony(ctx,q,arg.ToAccountID,arg.Amount,arg.FromAccountID,-arg.Amount)
		}
		if err!=nil{
			return err
		}
		return nil
	})

	return result,err
}

func addMony(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
)(account1 Account,account2 Account,err error){
	account1,err = q.AddAccountBalance(ctx,AddAccountBalanceParams{
		ID: accountID1,
		Amount: amount1,
	})
	if err!=nil{
		return
	}
	account2,err = q.AddAccountBalance(ctx,AddAccountBalanceParams{
		ID: accountID2,
		Amount: amount2,
	})
	return
}
	

