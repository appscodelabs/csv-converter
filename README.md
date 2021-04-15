# csv-converter

csv-converter \
  --in=/home/tamal/Downloads/mailchimp/cleaned_members_export_0e633f6c70.csv \
  --out=/home/tamal/Downloads/mailchimp \
  --renames confirm_ip=ip

csv-converter \
  --in=/home/tamal/Downloads/mailchimp/subscribed_members_export_0e633f6c70.csv \
  --out=/home/tamal/Downloads/mailchimp \
  --renames confirm_ip=ip

csv-converter \
  --in=/home/tamal/Downloads/mailchimp/unsubscribed_members_export_0e633f6c70.csv \
  --out=/home/tamal/Downloads/mailchimp \
  --renames confirm_ip=ip
