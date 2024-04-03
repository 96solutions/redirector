CREATE TABLE tracking_links
(
    slug                        varchar(55)               not null
        primary key,
    campaign_id                 int unsigned               not null,
    publisher_id                int unsigned               not null,
    advertiser_id               int unsigned               not null,
    source_id                   int unsigned               not null,
    active                      tinyint(1)   default 1     not null,
    allowed_protocols           json                       not null,
    allowed_geos                json                       not null,
    allowed_devices             json                       not null,

    campaign_overaged           tinyint(1)   default 0     not null,
    campaign_overaged_redirect_rules_id     int unsigned    not null,

    campaign_active           tinyint(1)   default 0     not null,
    campaign_active_redirect_rules_id     int unsigned    not null,

    target_url_template         varchar(255)            not null,

    created_at                  timestamp                  null,
    updated_at                  timestamp                  null
)
    collate = utf8_unicode_ci;
