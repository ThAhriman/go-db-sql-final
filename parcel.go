package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	defer db.Close()

	res, err := db.Exec("INSERT INTO parcel (client, status, address, created_At) VALUES (:client, :status, :address, :created_At)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_At", p.Created_At))
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	// верните идентификатор последней добавленной записи
	resId, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return int(resId), err
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number

	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		fmt.Println(err)
		return Parcel{}, err
	}
	defer db.Close()

	row := db.QueryRow("SELECT * FROM parcel WHERE number = :number",
		sql.Named("number", number))

	// здесь из таблицы должна вернуться только одна строка

	// заполните объект Parcel данными из таблицы

	p := Parcel{}
	err = row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.Created_At)
	if err != nil {
		fmt.Println(err)
		return Parcel{}, err
	}

	return p, nil

}
func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк

	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM parcel WHERE client = :client", sql.Named("client", client))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	// заполните срез Parcel данными из таблицы
	var res []Parcel

	for rows.Next() {

		cl := Parcel{}
		err := rows.Scan(&cl.Number, &cl.Client, &cl.Status, &cl.Address, &cl.Created_At)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		res = append(res, cl)

	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE parcel SET address = :address WHERE number = :number AND status = :status",
		sql.Named("number", number),
		sql.Named("address", address),
		sql.Named("status", ParcelStatusRegistered))

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM parcel WHERE number = :number AND status = :status",
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
