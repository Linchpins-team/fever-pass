use tcfsh;
/* recorded */
select *
from accounts a 
inner join (
    select max(id) as id, account_id
    from records
    where deleted_at is NULL
    group by account_id
) r
on a.id = r.account_id;

/* unrecorded */
select *
from accounts a 
left join (
    select max(id) as id, account_id
    from records
    where deleted_at is NULL
    group by account_id
) r
on a.id = r.account_id
where r.id is NULL;

# fevered 
select *
from accounts a 
inner join (
    select max(id) as id, account_id
    from records
    where deleted_at is NULL
    group by account_id
) m
inner join records r
on m.id = r.id
on a.id = r.account_id
where (temperature > 38 and type = 1) or (temperature > 37.5 and type = 2);