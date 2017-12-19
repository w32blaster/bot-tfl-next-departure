package db

import (
	"encoding/binary"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/w32blaster/bot-tfl-next-departure/structs"
)

const (

	// StateStationID state when user selected a first station on the last step and now selected the second one
	StateStationID = "station"

	// StateBookmark state when user wanted to save a bookmark and now it enters its name
	StateBookmark = "bookmark"

	// MaxLengthBookmarkName limit name for a bookmark name
	MaxLengthBookmarkName = 50
)

var (
	// bucket names
	bucketUserToBookmarks = []byte("user-to-bookmarks")
	bucketBookmarks       = []byte("bookmarks")
	bucketState           = []byte("state")
)

// GetStateFor simply get the last state for a user
func GetStateFor(userID int) (*structs.State, error) {
	db, err := connect()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer db.Close()

	var bytesJSON []byte
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketState)
		bytesJSON = b.Get(itob(userID))
		return nil
	})

	if bytesJSON == nil {
		return nil, nil
	}

	var state structs.State
	json.Unmarshal(bytesJSON, &state)
	return &state, nil
}

// SaveStateForStationID save state to the database
func SaveStateForStationID(userID int, stateStationID string) error {

	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {

		// open bucket
		b := tx.Bucket(bucketState)

		state := structs.State{
			Command:   StateStationID,
			StationID: stateStationID,
		}

		bytesJSON, _ := json.Marshal(state)

		// Persist bytes to state bucket.
		return b.Put(itob(userID), bytesJSON)
	})
}

// SaveStateForBookmark save state to the database when user wants to save a bookmark
func SaveStateForBookmark(userID int, jouney *structs.JourneyRequest) error {

	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {

		// open bucket
		b := tx.Bucket(bucketState)

		state := structs.State{
			Command:        StateBookmark,
			JourneyRequest: *jouney,
		}

		bytesJSON, _ := json.Marshal(state)

		// Persist bytes to state bucket.
		return b.Put(itob(userID), bytesJSON)
	})
}

// DeleteStateFor clear state
func DeleteStateFor(userID int) error {

	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketState)
		return b.Delete(itob(userID))
	})
}

// SaveBookmark Save the bookmark for a given user
func SaveBookmark(userID int, bookmarkName string, journey *structs.JourneyRequest) error {
	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// if the name is too long, then shorten it
	if len(bookmarkName) > MaxLengthBookmarkName {
		bookmarkName = bookmarkName[0:MaxLengthBookmarkName] + "..."
	}

	return db.Update(func(tx *bolt.Tx) error {

		// open bucket
		b := tx.Bucket(bucketBookmarks)

		id, _ := b.NextSequence()

		bookmark := structs.Bookmark{
			Name:    bookmarkName,
			Journey: *journey,
		}

		bytesJSON, _ := json.Marshal(bookmark)

		// Persist bytes to state bucket.
		b.Put(itob(int(id)), bytesJSON)

		// then save this bookmark to user
		b = tx.Bucket(bucketUserToBookmarks)

		bytesUserID := itob(userID)
		userBookmarksIDs := b.Get(bytesUserID)
		if userBookmarksIDs == nil {

			// create new relation
			b.Put(bytesUserID, []byte(strconv.Itoa(int(id))))

		} else {

			// update existing relationship
			strListOfIDs := string(userBookmarksIDs) + "," + strconv.Itoa(int(id))
			b.Put(bytesUserID, []byte(strListOfIDs))
		}

		return nil
	})
}

// GetBookmarksFor returns all the bookmakrs for an user
func GetBookmarksFor(userID int) *[]structs.Bookmark {

	db, err := connect()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer db.Close()

	var arrBookmarks []structs.Bookmark
	db.View(func(tx *bolt.Tx) error {

		// 1. Get list of IDs
		b := tx.Bucket(bucketUserToBookmarks)
		bytesList := b.Get(itob(userID))

		if bytesList == nil {
			return nil
		}
		// we store list as a string with IDs separated with comma
		arrIds := strings.Split(string(bytesList), ",")
		if len(arrIds) > 0 {
			arrBookmarks = make([]structs.Bookmark, len(arrIds))
			for i, bookmakrID := range arrIds {
				intID, _ := strconv.Atoi(bookmakrID)
				bookmark := getBookmark(tx, itob(intID))
				if bookmark != nil {
					arrBookmarks[i] = *bookmark
				}
			}
		}

		return nil
	})

	return &arrBookmarks
}

// gets one bookmark within one transaction
func getBookmark(tx *bolt.Tx, ID []byte) *structs.Bookmark {

	b := tx.Bucket(bucketBookmarks)
	bytesJSON := b.Get(ID)

	if bytesJSON == nil {
		return nil
	}

	var bookmark structs.Bookmark
	json.Unmarshal(bytesJSON, &bookmark)
	return &bookmark
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// Init initiates the database, creates all the buckets
func Init() error {

	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	return db.Batch(func(tx *bolt.Tx) error {

		_, err := tx.CreateBucketIfNotExists(bucketUserToBookmarks)
		if err != nil {
			log.Fatal(err)
			return err
		}

		_, err = tx.CreateBucketIfNotExists(bucketState)
		if err != nil {
			log.Fatal(err)
			return err
		}

		_, err = tx.CreateBucketIfNotExists(bucketBookmarks)
		if err != nil {
			log.Fatal(err)
			return err
		}

		return nil
	})
}

func connect() (*bolt.DB, error) {
	return bolt.Open("bot.db", 0600, &bolt.Options{Timeout: 10 * time.Second})
}
