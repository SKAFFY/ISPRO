package query

import (
	"github.com/Masterminds/squirrel"
)

var sb = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func ListLinks() squirrel.SelectBuilder {
	return sb.Select("id, source_id, target_id, created_at").From("links").OrderBy("id")
}

func GetLink(id int64) squirrel.SelectBuilder {
	return sb.Select("id, source_id, target_id, created_at").From("links").Where(squirrel.Eq{"id": id})
}

func CreateLink(sourceID, targetID int64) squirrel.InsertBuilder {
	return sb.Insert("links").
		Columns("source_id", "target_id").
		Values(sourceID, targetID).
		Suffix("RETURNING id, source_id, target_id, created_at")
}

func UpdateLink(id int64, sourceID, targetID int64) squirrel.UpdateBuilder {
	return sb.Update("links").
		Set("source_id", sourceID).
		Set("target_id", targetID).
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING id, source_id, target_id, created_at")
}

func DeleteLink(id int64) squirrel.DeleteBuilder {
	return sb.Delete("links").Where(squirrel.Eq{"id": id})
}
