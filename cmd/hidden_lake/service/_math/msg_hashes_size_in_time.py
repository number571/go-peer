def get_sizes(b):
    return f"\n\t{int(b)}B,\n\t{b/1024}KiB,\n\t{b/1024/1024}MiB,\n\t{b/1024/1024/1024}GiB"

HASH_SIZE  = 64 # hash(sha256) of message + hash(sha256) of public key
NODE_COUNT = 10 # N
PERIOD     = 5 # seconds

size_in_second = (HASH_SIZE * NODE_COUNT) / PERIOD
size_in_minute = size_in_second * 60
size_in_hour   = size_in_minute * 60
size_in_day    = size_in_hour * 24
size_in_week   = size_in_day * 7
size_in_month  = size_in_week * 4 + 2 # 4 weeks != full month
size_in_year   = size_in_month * 12

print("Second\t=", get_sizes(size_in_second))
print("Minute\t=", get_sizes(size_in_minute))
print("Hour\t=", get_sizes(size_in_hour))
print("Day\t=", get_sizes(size_in_day))
print("Week\t=", get_sizes(size_in_week))
print("Month\t=", get_sizes(size_in_month))
print("Year\t=", get_sizes(size_in_year))
