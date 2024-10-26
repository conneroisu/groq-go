echo "APIKEY: $COMPOSIO_API_KEY"
apikey=$COMPOSIO_API_KEY
curl --request GET --url https://backend.composio.dev/api/v2/actions?tags=Authentication --header 'X-API-Key: '$apikey \
> composio.json
