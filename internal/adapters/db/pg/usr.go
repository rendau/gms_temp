package pg

import (
	"context"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/rendau/gms_temp/internal/domain/entities"
)

func (d *St) UsrList(ctx context.Context, pars *entities.UsrListParsSt) ([]*entities.UsrListSt, int64, error) {
	var err error

	qFrom := ` from usr u`
	qWhere := ` where 1=1`
	qOffset := ``
	qLimit := ``
	args := map[string]interface{}{}

	// filter
	if pars.Ids != nil {
		args["ids"] = *pars.Ids
		qWhere += ` and u.id in (select * from unnest(${ids} :: bigint[]))`
	}

	if pars.TypeId != nil {
		args["type_id"] = *pars.TypeId
		qWhere += ` and u.type_id = ${type_id}`
	}

	if pars.TypeIds != nil {
		args["type_ids"] = *pars.TypeIds
		qWhere += ` and u.type_id in (select * from unnest(${type_ids} :: bigint[]))`
	}

	if pars.Search != nil && *pars.Search != "" {
		for wordI, word := range strings.Split(*pars.Search, " ") {
			if word != "" {
				key := "s_word" + strconv.Itoa(wordI)

				args[key] = word
				qWhere += ` and (
					u.phone ilike '%'|| ${` + key + `} ||'%' or
					u.name ilike '%'|| ${` + key + `} ||'%'
				)`
			}
		}
	}

	var tCount int64

	if pars.PageSize > 0 || pars.OnlyCount {
		err = d.DbQueryRowM(ctx, `select count(*)`+qFrom+qWhere, args).Scan(&tCount)
		if err != nil {
			return nil, 0, d.handleError(ctx, err)
		}

		if pars.OnlyCount {
			return nil, tCount, nil
		}

		qOffset = ` offset ` + strconv.FormatInt(pars.Page*pars.PageSize, 10)
		qLimit = ` limit ` + strconv.FormatInt(pars.PageSize, 10)
	}

	qSelect := `
		select
			  u.id
			, u.created_at
			, u.type_id
			, u.phone
			, u.ava
			, u.name
	`

	qOrderBy := ` order by u.name collate "C", u.id`
	if pars.SortBy != nil {
		qOrderBy = ` order by true`

		switch *pars.SortBy {
		case "phone":
			qOrderBy += ` , u.phone`
		case "name":
			qOrderBy += ` , u.name`
		}
	}

	rows, err := d.DbQueryM(ctx, qSelect+qFrom+qWhere+qOrderBy+qOffset+qLimit, args)
	if err != nil {
		return nil, 0, d.handleError(ctx, err)
	}
	defer rows.Close()

	recs := make([]*entities.UsrListSt, 0)

	for rows.Next() {
		rec := &entities.UsrListSt{}
		err = rows.Scan(
			&rec.Id,
			&rec.CreatedAt,
			&rec.TypeId,
			&rec.Phone,
			&rec.Ava,
			&rec.Name,
		)
		if err != nil {
			return nil, 0, d.handleError(ctx, err)
		}
		recs = append(recs, rec)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, d.handleError(ctx, err)
	}

	return recs, tCount, nil
}

func (d *St) UsrGet(ctx context.Context, pars *entities.UsrGetParsSt) (*entities.UsrSt, error) {
	args := make(map[string]interface{})

	qFrom := `usr u`
	qWhere := `1=1`

	if pars.Id != nil {
		args["id"] = *pars.Id
		qWhere += ` and u.id = ${id}`
	}

	if pars.Phone != nil {
		args["phone"] = *pars.Phone
		qWhere += ` and u.phone = ${phone}`
	}

	if pars.Token != nil {
		args["token"] = *pars.Token
		qWhere += ` and u.token = ${token}`
	}

	usr := &entities.UsrSt{}

	err := d.DbQueryRowM(ctx, `
		select
			  u.id
			, u.created_at
			, u.type_id
			, u.phone
			, u.ava
			, u.name
		from `+qFrom+`
		where `+qWhere+`
		limit 1
	`, args).Scan(
		&usr.Id,
		&usr.CreatedAt,
		&usr.TypeId,
		&usr.Phone,
		&usr.Ava,
		&usr.Name,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, d.handleError(ctx, err)
	}

	return usr, nil
}

func (d *St) UsrIdExists(ctx context.Context, id int64) (bool, error) {
	var cnt int

	err := d.DbQueryRow(ctx, `
		select count(*)
		from usr
		where id = $1
	`, id).Scan(&cnt)
	if err != nil {
		return false, d.handleError(ctx, err)
	}

	return cnt > 0, nil
}

func (d *St) UsrIdsExists(ctx context.Context, ids []int64) (bool, error) {
	var cnt int

	err := d.DbQueryRow(ctx, `
		select count(*)
		from unnest($1::bigint[]) x(id)
			left join usr u on u.id = x.id
		where u.id is null
	`, ids).Scan(&cnt)
	if err != nil {
		return false, d.handleError(ctx, err)
	}

	return cnt == 0, nil
}

func (d *St) UsrPhoneExists(ctx context.Context, phone string, excludeId int64) (bool, error) {
	var cnt int

	err := d.DbQueryRow(ctx, `
		select count(*)
		from usr
		where phone = $1
			and id != $2
	`, phone, excludeId).Scan(&cnt)
	if err != nil {
		return false, d.handleError(ctx, err)
	}

	return cnt > 0, nil
}

func (d *St) UsrGetToken(ctx context.Context, id int64) (string, error) {
	var token string

	err := d.DbQueryRow(ctx, `
		select token
		from usr
		where id = $1
	`, id).Scan(&token)
	if err != nil {
		return "", d.handleError(ctx, err)
	}

	return token, nil
}

func (d *St) UsrSetToken(ctx context.Context, id int64, token string) error {
	_, err := d.DbExec(ctx, `
		update usr set token = $2 where id = $1
	`, id, token)
	if err != nil {
		return d.handleError(ctx, err)
	}

	return nil
}

func (d *St) UsrGetTypeId(ctx context.Context, id int64) (int, error) {
	var result int

	err := d.DbQueryRow(ctx, `
		select type_id
		from usr
		where id = $1
	`, id).Scan(&result)
	if err != nil {
		return 0, d.handleError(ctx, err)
	}

	return result, nil
}

func (d *St) UsrGetPhone(ctx context.Context, id int64) (string, error) {
	var result string

	err := d.DbQueryRow(ctx, `
			select phone
			from usr
			where id = $1
	`, id).Scan(&result)
	if err != nil {
		return "", d.handleError(ctx, err)
	}

	return result, nil
}

func (d *St) UsrGetIdForPhone(ctx context.Context, phone string) (int64, error) {
	var result int64

	err := d.DbQueryRow(ctx, `
			select id
			from usr
			where phone = $1
	`, phone).Scan(&result)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}

		return 0, d.handleError(ctx, err)
	}

	return result, nil
}

func (d *St) UsrCreate(ctx context.Context, obj *entities.UsrCUSt) (int64, error) {
	args, err := d.UsrGetCUArgs(obj)
	if err != nil {
		return 0, err
	}

	var fields string
	var values string

	for k := range args {
		if fields != `` {
			fields += `,`
			values += `,`
		}
		fields += k
		values += `${` + k + `}`
	}

	var newId int64

	err = d.DbQueryRowM(ctx, `
		insert into usr(`+fields+`)
		values (`+values+`)
		returning id
	`, args).Scan(&newId)
	if err != nil {
		return 0, d.handleError(ctx, err)
	}

	return newId, nil
}

func (d *St) UsrUpdate(ctx context.Context, id int64, obj *entities.UsrCUSt) error {
	args, err := d.UsrGetCUArgs(obj)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return nil
	}

	var fields string

	for k := range args {
		if fields != `` {
			fields += `,`
		}
		fields += k + `=${` + k + `}`
	}

	args["id"] = id

	_, err = d.DbExecM(ctx, `
		update usr
		set `+fields+`
		where id = ${id}
	`, args)
	if err != nil {
		return d.handleError(ctx, err)
	}

	return nil
}

func (d *St) UsrGetCUArgs(obj *entities.UsrCUSt) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	if obj.TypeId != nil {
		result["type_id"] = *obj.TypeId
	}

	if obj.Phone != nil {
		result["phone"] = *obj.Phone
	}

	if obj.Ava != nil {
		result["ava"] = *obj.Ava
	}

	if obj.Name != nil {
		result["name"] = *obj.Name
	}

	return result, nil
}

func (d *St) UsrDelete(ctx context.Context, id int64) error {
	_, err := d.DbExec(ctx, `delete from usr where id = $1`, id)
	if err != nil {
		return d.handleError(ctx, err)
	}

	return nil
}

func (d *St) UsrFilterUnusedFiles(ctx context.Context, src []string) ([]string, error) {
	rows, err := d.DbQuery(ctx, `
		select x.a
		from unnest($1 :: text[]) x(a)
			left join usr y on y.ava = x.a
		where y.ava is null
	`, src)
	if err != nil {
		return nil, d.handleError(ctx, err)
	}
	defer rows.Close()

	result := make([]string, 0)

	var v string

	for rows.Next() {
		err = rows.Scan(&v)
		if err != nil {
			return nil, d.handleError(ctx, err)
		}

		result = append(result, v)
	}
	if err = rows.Err(); err != nil {
		return nil, d.handleError(ctx, err)
	}

	return result, nil
}
