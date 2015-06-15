package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Connection struct {
	db *sql.DB
}

type JobNetwork struct {
	ID     int
	Name   string
	Start  string
	End    string
	Status int
}

type Job struct {
	NID    int
	JID    string
	Name   string
	Start  string
	End    string
	Status int
	RC     int
}

const driver = "sqlite3"

func Open(dbfile string) (*Connection, error) {
	if _, exist := os.Stat(dbfile); exist != nil {
		return nil, fmt.Errorf("Not found dbfile[%v]", dbfile)
	}
	db, err := sql.Open(driver, dbfile)
	if err != nil {
		return nil, err
	}

	return &Connection{db}, nil
}

func (c *Connection) SelectJobNetwork(id int) (*JobNetwork, error) {
	query := "SELECT JOBNETWORK, STARTDATE, ENDDATE, STATUS FROM JOBNETWORK WHERE ID = ?"
	row := c.db.QueryRow(query, id)
	n := new(JobNetwork)
	err := row.Scan(&n.Name, &n.Start, &n.End, &n.Status)
	if err != nil {
		return nil, err
	}
	n.ID = id
	return n, nil
}

func (c *Connection) SelectJob(nid int, jid string) (*Job, error) {
	query := "SELECT JOBNAME, STARTDATE, ENDDATE, STATUS, RC FROM JOB WHERE ID = ? AND JOBID = ?"
	row := c.db.QueryRow(query, nid, jid)
	j := new(Job)
	err := row.Scan(&j.Name, &j.Start, &j.End, &j.Status, &j.RC)
	if err != nil {
		return nil, err
	}
	j.NID = nid
	j.JID = jid
	return j, nil
}

func (c *Connection) Close() {
	c.db.Close()
}
