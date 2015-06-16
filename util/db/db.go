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

func (c *Connection) SelectJobNetworksByCond(cond string) ([]*JobNetwork, error) {
	query := "SELECT ID, JOBNETWORK, STARTDATE, ENDDATE, STATUS FROM JOBNETWORK WHERE " + cond
	networks := make([]*JobNetwork, 0)
	rows, err := c.db.Query(query)
	if err != nil {
		return networks, err
	}

	for rows.Next() {
		n := new(JobNetwork)
		rows.Scan(&n.ID, &n.Name, &n.Start, &n.End, &n.Status)
		networks = append(networks, n)
	}
	return networks, nil
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

func (c *Connection) SelectJobsByCond(condition string) ([]*Job, error) {
	query := "SELECT ID, JOBID, JOBNAME, STARTDATE, ENDDATE, STATUS, RC FROM JOB WHERE " + condition
	jobs := make([]*Job, 0)
	rows, err := c.db.Query(query)
	if err != nil {
		return jobs, err
	}
	defer rows.Close()

	for rows.Next() {
		j := new(Job)
		rows.Scan(&j.NID, &j.JID, &j.Name, &j.Start, &j.End, &j.Status, &j.RC)
		jobs = append(jobs, j)
	}
	return jobs, nil
}

func (c *Connection) CountJobs(nid int) (int, error) {
	query := "SELECT COUNT(*) FROM JOB WHERE ID = ?"
	row := c.db.QueryRow(query, nid)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

func (c *Connection) Close() {
	c.db.Close()
}
