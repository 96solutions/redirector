CREATE TABLE tracking_links
(
    slug                        varchar(55)               not null
        primary key,
    campaign_id                 int                       not null,
    affiliate_id                int                       not null,
    advertiser_id               int                       not null,
    source_id                   int                       not null,
    active                      boolean    default true   not null,
    allowed_protocols           jsonb                     not null,
    allowed_geos                jsonb                     not null,
    allowed_devices             jsonb                     not null,

    campaign_overaged           boolean    default false  not null,
    campaign_overaged_redirect_rules_id     int           not null,

    campaign_active             boolean    default false  not null,
    campaign_active_redirect_rules_id       int           not null,

    target_url_template         varchar(255)             not null,

    created_at                  timestamp without time zone null,
    updated_at                  timestamp without time zone null
);
