# csv-converter

csv-to-json \
  --in=/home/tamal/Downloads/mailchimp/cleaned_members_export_0e633f6c70.csv \
  --out=/home/tamal/Downloads/mailchimp \
  --renames confirm_ip=ip

csv-to-json \
  --in=/home/tamal/Downloads/mailchimp/subscribed_members_export_0e633f6c70.csv \
  --out=/home/tamal/Downloads/mailchimp \
  --renames confirm_ip=ip

csv-to-json \
  --in=/home/tamal/Downloads/mailchimp/unsubscribed_members_export_0e633f6c70.csv \
  --out=/home/tamal/Downloads/mailchimp \
  --renames confirm_ip=ip

filter-json \
  --in=/home/tamal/Downloads/mailchimp/subscribed_members_export_0e633f6c70.json \
  --keys=cc \
  --keys=email \
  --keys=ip \
  --keys=latitude \
  --keys=longitude \
  --keys=timezone

json-to-listmonk-csv \
  --in=/home/tamal/Downloads/mailchimp/subscribed_members_export_0e633f6c70_filtered.json
