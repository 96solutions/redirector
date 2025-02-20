-- Add new columns for redirect rules
ALTER TABLE tracking_links
    ADD COLUMN campaign_protocol_redirect_rules_id integer REFERENCES redirect_rules(id),
    ADD COLUMN campaign_geo_redirect_rules_id integer REFERENCES redirect_rules(id),
    ADD COLUMN campaign_devices_redirect_rules_id integer REFERENCES redirect_rules(id),
    ADD COLUMN campaign_os_redirect_rules_id integer REFERENCES redirect_rules(id),
    ADD COLUMN allowed_os jsonb NOT NULL DEFAULT '{}',
    ADD COLUMN allow_deeplink boolean NOT NULL DEFAULT false;

-- Add foreign key constraints for existing redirect rules columns
ALTER TABLE tracking_links
    ADD FOREIGN KEY (campaign_overaged_redirect_rules_id) REFERENCES redirect_rules(id),
    ADD FOREIGN KEY (campaign_active_redirect_rules_id) REFERENCES redirect_rules(id);

-- Update ID columns to be consistent with entity type
ALTER TABLE tracking_links
    ALTER COLUMN campaign_id TYPE varchar(255),
    ALTER COLUMN affiliate_id TYPE varchar(255),
    ALTER COLUMN advertiser_id TYPE varchar(255),
    ALTER COLUMN source_id TYPE varchar(255);
