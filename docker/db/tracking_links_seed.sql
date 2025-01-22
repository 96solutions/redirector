INSERT INTO tracking_links (
    slug, campaign_id, publisher_id, advertiser_id, source_id, active,
    allowed_protocols, allowed_geos, allowed_devices,
    campaign_overaged, campaign_overaged_redirect_rules_id,
    campaign_active, campaign_active_redirect_rules_id,
    target_url_template, created_at, updated_at
) VALUES (
             'slug-001', 101, 201, 301, 401, true,
             '{"http": true, "https": true}', '{"US": true, "UK": true}', '{"mobile": true, "desktop": true}',
             false, 1001,
             true, 2001,
             'https://httpbin.org/get?ip={ip}&click_id={click_id}&user_agent={user_agent}&campaign_id={campaign_id}&aff_id={aff_id}&source_id={source_id}&advertiser_id={advertiser_id}&date={date}&date_time={date_time}&timestamp={timestamp}&p1={p1}&p2={p2}&p3={p3}&p4={p4}&country_code={country_code}&referer={referer}&random_str={random_str}&random_int={random_int}&device={device}&platform={platform}', NOW(), NOW()
         );

INSERT INTO tracking_links (
    slug, campaign_id, publisher_id, advertiser_id, source_id, active,
    allowed_protocols, allowed_geos, allowed_devices,
    campaign_overaged, campaign_overaged_redirect_rules_id,
    campaign_active, campaign_active_redirect_rules_id,
    target_url_template, created_at, updated_at
) VALUES (
             'slug-002', 102, 202, 302, 402, true,
             '{"http": false, "https": true}', '{"CA": true, "AU": true}', '{"tablet": true, "desktop": true}',
             true, 1002,
             false, 2002,
             'https://httpbin.org/get?ip={ip}&click_id={click_id}&user_agent={user_agent}&campaign_id={campaign_id}&aff_id={aff_id}&source_id={source_id}&advertiser_id={advertiser_id}&date={date}&date_time={date_time}&timestamp={timestamp}&p1={p1}&p2={p2}&p3={p3}&p4={p4}&country_code={country_code}&referer={referer}&random_str={random_str}&random_int={random_int}&device={device}&platform={platform}', NOW(), NOW()
         );

INSERT INTO tracking_links (
    slug, campaign_id, publisher_id, advertiser_id, source_id, active,
    allowed_protocols, allowed_geos, allowed_devices,
    campaign_overaged, campaign_overaged_redirect_rules_id,
    campaign_active, campaign_active_redirect_rules_id,
    target_url_template, created_at, updated_at
) VALUES (
             'slug-003', 103, 203, 303, 403, false,
             '{"http": true, "https": true}', '{"DE": true, "FR": true}', '{"mobile": true, "tablet": true}',
             false, 1003,
             true, 2003,
             'https://httpbin.org/get?ip={ip}&click_id={click_id}&user_agent={user_agent}&campaign_id={campaign_id}&aff_id={aff_id}&source_id={source_id}&advertiser_id={advertiser_id}&date={date}&date_time={date_time}&timestamp={timestamp}&p1={p1}&p2={p2}&p3={p3}&p4={p4}&country_code={country_code}&referer={referer}&random_str={random_str}&random_int={random_int}&device={device}&platform={platform}', NOW(), NOW()
         );

INSERT INTO tracking_links (
    slug, campaign_id, publisher_id, advertiser_id, source_id, active,
    allowed_protocols, allowed_geos, allowed_devices,
    campaign_overaged, campaign_overaged_redirect_rules_id,
    campaign_active, campaign_active_redirect_rules_id,
    target_url_template, created_at, updated_at
) VALUES (
             'slug-004', 104, 204, 304, 404, true,
             '{"http": false, "https": true}', '{"IN": true, "BR": true, "PL": true}', '{"mobile": true}',
             true, 1004,
             false, 2004,
             'https://httpbin.org/get?ip={ip}&click_id={click_id}&user_agent={user_agent}&campaign_id={campaign_id}&aff_id={aff_id}&source_id={source_id}&advertiser_id={advertiser_id}&date={date}&date_time={date_time}&timestamp={timestamp}&p1={p1}&p2={p2}&p3={p3}&p4={p4}&country_code={country_code}&referer={referer}&random_str={random_str}&random_int={random_int}&device={device}&platform={platform}', NOW(), NOW()
         );

INSERT INTO tracking_links (
    slug, campaign_id, publisher_id, advertiser_id, source_id, active,
    allowed_protocols, allowed_geos, allowed_devices,
    campaign_overaged, campaign_overaged_redirect_rules_id,
    campaign_active, campaign_active_redirect_rules_id,
    target_url_template, created_at, updated_at
) VALUES (
             'slug-005', 105, 205, 305, 405, false,
             '{"http": true, "https": true}', '{"JP": true, "KR": true}', '{"desktop": true}',
             false, 1005,
             true, 2005,
             'https://httpbin.org/get?ip={ip}&click_id={click_id}&user_agent={user_agent}&campaign_id={campaign_id}&aff_id={aff_id}&source_id={source_id}&advertiser_id={advertiser_id}&date={date}&date_time={date_time}&timestamp={timestamp}&p1={p1}&p2={p2}&p3={p3}&p4={p4}&country_code={country_code}&referer={referer}&random_str={random_str}&random_int={random_int}&device={device}&platform={platform}', NOW(), NOW()
         );
