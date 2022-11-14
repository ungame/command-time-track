package repository

const (
	insertActivityQuery = `insert into activities (category, description, started_at, updated_at) values (?, ?, ?, ?)`
	updateActivityQuery = `update activities set category = ?, description = ?, updated_at = ?, stopped_at = ? where id = ?`
	deleteActivityQuery = `delete from activities where id = ?`
)
