package db

import (
	"context"
	"database/sql"
	"fmt"
)

//generic interface for all the queries
//anyone who wants to be a store must implement this interface
type Store interface {
	Querier //genrated by sqlc
	LikeTx(ctx context.Context, arg LikePostParams) error
}

// has db and set of queries to interact with the database
type SQLStore struct {
	db *sql.DB
	*Queries
}

// constructor
func NewStore(db *sql.DB) Store {

	//sqlstore implements the store interface as the queries field is a part of it
	//and the likeTx method is implemented by the sqlstore
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// creates a closure that executes a function within a database transaction
// whatever queries we write in the @fn will be executed in a transaction
func (store SQLStore) execTransaction(ctx context.Context, fn func(*Queries) error) error {

	//create a transaction
	transaction, err := store.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	//we pass in the trasaction object to new as it is an interface that implements the both db and tx
	q := New(transaction)

	//execute the function
	err = fn(q)
	if err != nil {
		//if there is an error, rollback the transaction
		if rbErr := transaction.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	//if there is no error, commit the transaction
	return transaction.Commit()
}

type LikePostParams struct {
	UserID int32 `json:"user_id"`
	PostID int32 `json:"post_id"`
}

// like a post and update the reputation of the author
func (store *SQLStore) LikeTx(ctx context.Context, arg LikePostParams) error {
	err := store.execTransaction(ctx, func(q *Queries) error {
		var err error
		//like the post
		err = q.LikePost(ctx, int64(arg.UserID))
		if err != nil {
			return err
		}

		//TODO: solve deadlock
		//update the reputation of the author
		err = q.UpdateReputation(ctx, UpdateReputationParams{
			ID:         int64(arg.UserID),
			Reputation: 10,
		})
		return err
	})

	//get the new like and reputation count
	return err
}
