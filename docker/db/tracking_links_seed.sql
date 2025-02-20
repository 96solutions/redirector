-- Insert redirect rules
INSERT INTO redirect_rules (id, redirect_type, redirect_slug, redirect_url, redirect_smart_slug)
VALUES (1001, 'url', NULL, 'https://overaged1.example.com', NULL),
       (1002, 'url', NULL, 'https://overaged2.example.com', NULL),
       (1003, 'smart_slug', NULL, NULL, '[
         "slug-001",
         "slug-002",
         "slug-003"
       ]'),
       (1004, 'slug', 'alternative-slug', NULL, NULL),
       (1005, 'url', NULL, 'https://overaged5.example.com', NULL),
       (2001, 'url', NULL, 'https://disabled1.example.com', NULL),
       (2002, 'smart_slug', NULL, NULL, '[
         "slug-001",
         "slug-002",
         "slug-003"
       ]'),
       (2003, 'slug', 'disabled-slug', NULL, NULL),
       (2004, 'url', NULL, 'https://disabled4.example.com', NULL),
       (2005, 'url', NULL, 'https://disabled5.example.com', NULL),
       (3001, 'url', NULL, 'https://protocol1.example.com', NULL),
       (3002, 'url', NULL, 'https://protocol2.example.com', NULL),
       (3003, 'url', NULL, 'https://protocol3.example.com', NULL),
       (4001, 'url', NULL, 'https://geo1.example.com', NULL),
       (4002, 'url', NULL, 'https://geo2.example.com', NULL),
       (5001, 'url', NULL, 'https://device1.example.com', NULL),
       (5002, 'url', NULL, 'https://device2.example.com', NULL),
       (6001, 'url', NULL, 'https://os1.example.com', NULL),
       (6002, 'url', NULL, 'https://os2.example.com', NULL);

-- Insert tracking links with updated fields
INSERT INTO tracking_links (slug, campaign_id, affiliate_id, advertiser_id, source_id, active,
                            allowed_protocols, allowed_geos, allowed_devices, allowed_os,
                            campaign_overaged, campaign_overaged_redirect_rules_id,
                            campaign_active, campaign_active_redirect_rules_id,
                            campaign_protocol_redirect_rules_id,
                            campaign_geo_redirect_rules_id,
                            campaign_devices_redirect_rules_id,
                            campaign_os_redirect_rules_id,
                            target_url_template,
                            allow_deeplink,
                            created_at, updated_at)
VALUES ('slug-001',
        'campaign_101',
        'affiliate_201',
        'advertiser_301',
        'source_401',
        true,
        '{
          "http": true,
          "https": true
        }',
        '{
          "US": true,
          "UK": true
        }',
        '{
          "mobile": true,
          "desktop": true
        }',
        '{
          "android": true,
          "ios": true
        }',
        false,
        1001,
        true,
        2001,
        3001,
        4001,
        5001,
        6001,
        'https://httpbin.org/get?ip={ip}&click_id={click_id}&user_agent={user_agent}&campaign_id={campaign_id}&aff_id={aff_id}&source_id={source_id}&advertiser_id={advertiser_id}&date={date}&date_time={date_time}&timestamp={timestamp}&p1={p1}&p2={p2}&p3={p3}&p4={p4}&country_code={country_code}&referer={referer}&random_str={random_str}&random_int={random_int}&device={device}&platform={platform}',
        true,
        NOW(),
        NOW());

INSERT INTO tracking_links (slug, campaign_id, affiliate_id, advertiser_id, source_id, active,
                            allowed_protocols, allowed_geos, allowed_devices, allowed_os,
                            campaign_overaged, campaign_overaged_redirect_rules_id,
                            campaign_active, campaign_active_redirect_rules_id,
                            campaign_protocol_redirect_rules_id,
                            campaign_geo_redirect_rules_id,
                            campaign_devices_redirect_rules_id,
                            campaign_os_redirect_rules_id,
                            target_url_template,
                            allow_deeplink,
                            created_at, updated_at)
VALUES ('slug-002',
        'campaign_102',
        'affiliate_202',
        'advertiser_302',
        'source_402',
        true,
        '{
          "http": false,
          "https": true
        }',
        '{
          "CA": true,
          "AU": true
        }',
        '{
          "tablet": true,
          "desktop": true
        }',
        '{
          "windows": true,
          "macos": true
        }',
        true,
        1002,
        false,
        2002,
        3002,
        4002,
        5002,
        6002,
        'https://httpbin.org/get?ip={ip}&click_id={click_id}&user_agent={user_agent}&campaign_id={campaign_id}&aff_id={aff_id}&source_id={source_id}&advertiser_id={advertiser_id}&date={date}&date_time={date_time}&timestamp={timestamp}&p1={p1}&p2={p2}&p3={p3}&p4={p4}&country_code={country_code}&referer={referer}&random_str={random_str}&random_int={random_int}&device={device}&platform={platform}',
        false,
        NOW(),
        NOW());

INSERT INTO tracking_links (slug, campaign_id, affiliate_id, advertiser_id, source_id, active,
                            allowed_protocols, allowed_geos, allowed_devices,
                            campaign_overaged, campaign_overaged_redirect_rules_id,
                            campaign_active, campaign_active_redirect_rules_id,
                            target_url_template, created_at, updated_at)
VALUES ('slug-003', 103, 203, 303, 403, false,
        '{
          "http": true,
          "https": true
        }', '{
    "DE": true,
    "FR": true
  }', '{
    "mobile": true,
    "tablet": true
  }',
        false, 1003,
        true, 2003,
        'https://httpbin.org/get?ip={ip}&click_id={click_id}&user_agent={user_agent}&campaign_id={campaign_id}&aff_id={aff_id}&source_id={source_id}&advertiser_id={advertiser_id}&date={date}&date_time={date_time}&timestamp={timestamp}&p1={p1}&p2={p2}&p3={p3}&p4={p4}&country_code={country_code}&referer={referer}&random_str={random_str}&random_int={random_int}&device={device}&platform={platform}',
        NOW(), NOW());

INSERT INTO tracking_links (slug, campaign_id, affiliate_id, advertiser_id, source_id, active,
                            allowed_protocols, allowed_geos, allowed_devices,
                            campaign_overaged, campaign_overaged_redirect_rules_id,
                            campaign_active, campaign_active_redirect_rules_id,
                            target_url_template, created_at, updated_at)
VALUES ('slug-004', 104, 204, 304, 404, true,
        '{
          "http": false,
          "https": true
        }', '{
    "IN": true,
    "BR": true,
    "PL": true
  }', '{
    "mobile": true
  }',
        true, 1004,
        false, 2004,
        'https://httpbin.org/get?ip={ip}&click_id={click_id}&user_agent={user_agent}&campaign_id={campaign_id}&aff_id={aff_id}&source_id={source_id}&advertiser_id={advertiser_id}&date={date}&date_time={date_time}&timestamp={timestamp}&p1={p1}&p2={p2}&p3={p3}&p4={p4}&country_code={country_code}&referer={referer}&random_str={random_str}&random_int={random_int}&device={device}&platform={platform}',
        NOW(), NOW());

INSERT INTO tracking_links (slug, campaign_id, affiliate_id, advertiser_id, source_id, active,
                            allowed_protocols, allowed_geos, allowed_devices,
                            campaign_overaged, campaign_overaged_redirect_rules_id,
                            campaign_active, campaign_active_redirect_rules_id,
                            target_url_template, created_at, updated_at)
VALUES ('slug-005', 105, 205, 305, 405, false,
        '{
          "http": true,
          "https": true
        }', '{
    "JP": true,
    "KR": true
  }', '{
    "desktop": true
  }',
        false, 1005,
        true, 2005,
        'https://httpbin.org/get?ip={ip}&click_id={click_id}&user_agent={user_agent}&campaign_id={campaign_id}&aff_id={aff_id}&source_id={source_id}&advertiser_id={advertiser_id}&date={date}&date_time={date_time}&timestamp={timestamp}&p1={p1}&p2={p2}&p3={p3}&p4={p4}&country_code={country_code}&referer={referer}&random_str={random_str}&random_int={random_int}&device={device}&platform={platform}',
        NOW(), NOW());

-- Insert landing pages
INSERT INTO landing_pages (id, slug, title, preview_url, target_url)
VALUES ('landing_001', 'slug-001', 'Landing Page 1', 'https://preview1.example.com', 'https://target1.example.com'),
       ('landing_002', 'slug-001', 'Landing Page 2', 'https://preview2.example.com', 'https://target2.example.com'),
       ('landing_003', 'slug-002', 'Landing Page 3', 'https://preview3.example.com', 'https://target3.example.com');
