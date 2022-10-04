create table if not exists images (
    name text not null,
    location text,
	longitude float,
	latitude  float,
    primary key(name)
);

create or replace function isInsideBBox(swLong float, swLat float, neLong float, neLat float, pLong float, pLat float)
returns boolean language plpgsql as $$
declare
  isLongIn boolean;
begin
	if swLong < neLong then
		isLongIn := pLong >= swLong and pLong <= neLong;
	else
		isLongIn := pLong >= swLong or pLong <= neLong;
	end if;

	return pLat >= swLat and pLat <= neLat and isLongIn;
end $$;
