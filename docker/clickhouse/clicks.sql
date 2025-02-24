CREATE TABLE IF NOT EXISTS clicks (
                                      id String,
                                      target_url String,
                                      referer String,
                                      trk_url String,
                                      slug String,
                                      parent_slug String,

                                      source_id String,
                                      campaign_id String,
                                      affiliate_id String,
                                      advertiser_id String,
                                      is_parallel UInt8,

                                      landing_id String,
                                      gclid String,

                                      agent String,
                                      platform String,
                                      browser String,
                                      device String,

                                      ip String,
                                      country_code FixedString(2),

                                      p1 String,
                                      p2 String,
                                      p3 String,
                                      p4 String,

                                      created_at DateTime,

                                      INDEX idx_slug slug TYPE bloom_filter GRANULARITY 1,
                                      INDEX idx_campaign campaign_id TYPE bloom_filter GRANULARITY 1,
                                      INDEX idx_affiliate affiliate_id TYPE bloom_filter GRANULARITY 1,
                                      INDEX idx_source source_id TYPE bloom_filter GRANULARITY 1
)
    ENGINE = MergeTree()
PARTITION BY toYYYYMM(created_at)
ORDER BY (created_at, id)
SETTINGS index_granularity = 8192;