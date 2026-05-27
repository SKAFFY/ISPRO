package query

import (
	"github.com/Masterminds/squirrel"
)

var sb = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func ListEntries() squirrel.SelectBuilder {
	return sb.Select("id, title, content, created_at, updated_at").From("entries").OrderBy("id")
}

func GetEntry(id int64) squirrel.SelectBuilder {
	return sb.Select("id, title, content, created_at, updated_at").From("entries").Where(squirrel.Eq{"id": id})
}

func CreateEntry(title, content string) squirrel.InsertBuilder {
	return sb.Insert("entries").
		Columns("title", "content").
		Values(title, content).
		Suffix("RETURNING id, title, content, created_at, updated_at")
}

func UpdateEntry(id int64, title, content string) squirrel.UpdateBuilder {
	return sb.Update("entries").
		Set("title", title).
		Set("content", content).
		Set("updated_at", "NOW()").
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING id, title, content, created_at, updated_at")
}

func DeleteEntry(id int64) squirrel.DeleteBuilder {
	return sb.Delete("entries").Where(squirrel.Eq{"id": id})
}
