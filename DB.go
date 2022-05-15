package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type Storage interface {
	Add(album) (album, error)
	Update(string, album) (album, error)
	Delete(string) error
	Read() ([]album, error)
	OneRecord(string) (album, error)
}

type PostgreStorage struct {
	db *sql.DB
}

func (p PostgreStorage) CreateSchema() {
	_, err := p.db.Query("create table if not exists albums(" +
		"ID char(16) primary key," +
		"Title char(128)," +
		"Artist char(128)," +
		"Price decimal" +
		")")
	if err != nil {
		log.Fatal(err)
	}

}

func NewPostgresStorage() PostgreStorage {
	//docker run -it --name some-postgres -e POSTGRES_PASSWORD=pass -e POSTGRES_USER=user -e POSTGRES_DB=db postgres
	connStr := "user=user dbname=db password=pass sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	storage := PostgreStorage{db}
	storage.CreateSchema()
	return storage
}

func (p PostgreStorage) Add(al album) (album, error) {
	row := p.db.QueryRow("INSERT INTO albums(ID, Title, Artist, Price) "+
		"values($1, $2, $3, $4)", al.ID, al.Title, al.Artist, al.Price)
	if row.Err() != nil {
		return album{}, row.Err()
	}
	return al, nil
}

func (p PostgreStorage) OneRecord(id string) (album, error) {
	var alb album
	row := p.db.QueryRow("select * from albums where id = $1", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		return album{}, err
	}
	return alb, nil
}

func (p PostgreStorage) Update(id string, a album) (album, error) {
	row := p.db.QueryRow("update albums set Title = $1, Artist=$2, Price=$3 where id=$4", a.Title, a.Artist, a.Price, a.ID)
	if row.Err() != nil {
		return album{}, row.Err()
	}
	return a, nil
}

func (p PostgreStorage) Delete(id string) error {
	_, err := p.db.Query("delete from albums where id=$1", id)
	return err
}

func (p PostgreStorage) Read() ([]album, error) {
	var albums []album
	rows, err := p.db.Query("select * from albums")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var alb album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist,
			&alb.Price); err != nil {
			return albums, err
		}
		albums = append(albums, alb)
	}
	if err = rows.Err(); err != nil {
		return albums, err
	}
	return albums, nil
}

func NewStorage() Storage {
	return NewPostgresStorage()
}

//{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
//{ID: "2", Title: "Hypno", Artist: "Travis Scott", Price: 120.99},
//{ID: "3", Title: "Rich da kid", Artist: "Asap Rocky", Price: 573.99},
