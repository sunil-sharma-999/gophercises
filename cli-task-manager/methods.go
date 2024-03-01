package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

func (c *CLI) add(task string) error {

	err := c.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(c.Name))
		if err != nil {
			return err
		}

		id, err := bucket.NextSequence()

		if err != nil {
			return err
		}
		data, err := json.Marshal(Task{
			ID:        int(id),
			Text:      task,
			Done:      false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		if err != nil {
			return err
		}
		err = bucket.Put(Itob(int(id)), data)
		return err
	})

	return err
}

func (c *CLI) list(condition func(Task) bool) ([]Task, error) {
	results := []Task{}

	err := c.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(c.Name))
		if bucket == nil {
			fmt.Println("no tasks added yet")
			os.Exit(0)
			return nil
		}

		err := bucket.ForEach(func(k, v []byte) error {
			var data Task
			err := json.Unmarshal(v, &data)
			if err != nil {
				return nil
			}
			if condition == nil {
				results = append(results, data)

			} else {
				if condition(data) {
					results = append(results, data)
				}
			}
			return nil
		})

		return err

	})

	return results, err
}

func (c *CLI) do(id int) error {

	err := c.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(c.Name))
		if bucket == nil {
			fmt.Println("no tasks added yet")
			os.Exit(0)
			return nil
		}
		var dataJSON Task
		err := c.DB.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(c.Name))
			if bucket == nil {
				fmt.Println("no tasks added yet")
				os.Exit(0)
				return nil
			}
			data := bucket.Get(Itob(id))
			if data == nil {
				fmt.Printf("No task with ID: %v\n", id)
				os.Exit(0)
				return nil
			}
			err := json.Unmarshal(data, &dataJSON)
			return err
		})
		if err != nil {
			return err
		}
		if dataJSON.Done {
			fmt.Printf("You have already completed \"%v\" task.\n", dataJSON.Text)
			os.Exit(0)
			return nil
		}
		dataJSON.Done = true
		dataJSON.UpdatedAt = time.Now()
		marshalledData, err := json.Marshal(dataJSON)

		if err != nil {
			return err
		}
		err = bucket.Put(Itob(int(id)), marshalledData)
		if err != nil {
			return err
		}
		fmt.Printf("You have completed the \"%v\" task.\n", dataJSON.Text)
		return nil
	})

	return err
}

func (c *CLI) rm(id int) error {
	err := c.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(c.Name))
		if err != nil {
			return err
		}
		var dataJSON Task
		err = c.DB.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(c.Name))
			if bucket == nil {
				fmt.Println("no tasks added yet")
				os.Exit(0)
				return nil
			}
			data := bucket.Get(Itob(id))
			if data == nil {
				fmt.Printf("No task with ID: %v", id)
				os.Exit(0)
				return nil
			}
			err := json.Unmarshal(data, &dataJSON)
			return err
		})
		if err != nil {
			return err
		}

		err = bucket.Delete(Itob(int(id)))
		if err != nil {
			return err
		}
		fmt.Printf("You have deleted the \"%v\" task.\n", dataJSON.Text)
		return nil
	})

	return err
}
