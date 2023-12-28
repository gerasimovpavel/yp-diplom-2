create function gen_v7_uuid() returns uuid
    parallel safe
    language plpgsql
as
$$
declare
 tp text = lpad(to_hex(floor(extract(epoch from clock_timestamp())*1000)::int8),
 12,'0')||'7';
 entropy text = md5(gen_random_uuid()::text);
begin
 return (tp ||
 substring(entropy from 1 for 3) || to_hex(8+(random()*3)::int) ||
 substring(entropy from 4 for 15)
 )::uuid;
end
$$;

alter function gen_v7_uuid() owner to postgres;

