-- Remove foreign key constraints
ALTER TABLE tracking_links
    DROP CONSTRAINT IF EXISTS tracking_links_campaign_protocol_redirect_rules_id_fkey,
    DROP CONSTRAINT IF EXISTS tracking_links_campaign_geo_redirect_rules_id_fkey,
    DROP CONSTRAINT IF EXISTS tracking_links_campaign_devices_redirect_rules_id_fkey,
    DROP CONSTRAINT IF EXISTS tracking_links_campaign_os_redirect_rules_id_fkey,
    DROP CONSTRAINT IF EXISTS tracking_links_campaign_overaged_redirect_rules_id_fkey,
    DROP CONSTRAINT IF EXISTS tracking_links_campaign_active_redirect_rules_id_fkey;

-- Remove new columns
ALTER TABLE tracking_links
    DROP COLUMN IF EXISTS campaign_protocol_redirect_rules_id,
    DROP COLUMN IF EXISTS campaign_geo_redirect_rules_id,
    DROP COLUMN IF EXISTS campaign_devices_redirect_rules_id,
    DROP COLUMN IF EXISTS campaign_os_redirect_rules_id,
    DROP COLUMN IF EXISTS allowed_os,
    DROP COLUMN IF EXISTS allow_deeplink;

-- Revert ID columns to original types
ALTER TABLE tracking_links
    ALTER COLUMN campaign_id TYPE integer USING campaign_id::integer,
    ALTER COLUMN affiliate_id TYPE integer USING affiliate_id::integer,
    ALTER COLUMN advertiser_id TYPE integer USING advertiser_id::integer,
    ALTER COLUMN source_id TYPE integer USING source_id::integer;
