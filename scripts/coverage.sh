#!/bin/bash
#
# Code coverage generation

COVERAGE_DIR="${COVERAGE_DIR:-coverage}"

PACKAGE_LIST=$(go list ./... | grep -v -e '\/embed' -e '\/gen' -e '\/mocks')
COVER_PKG_LIST=$(echo "$PACKAGE_LIST" | tr "\n" ",")
PACKAGE_LIST=$(echo "$PACKAGE_LIST" | tr "\n" " ")

# Create the coverage files directory
mkdir -p "$COVERAGE_DIR"
rm -f "${COVERAGE_DIR}"/coverage.cov
rm -f "${COVERAGE_DIR}"/coverage.svg
rm -f "${COVERAGE_DIR}"/index.html

go test -cover \
 -covermode=atomic \
 -coverprofile "${COVERAGE_DIR}"/coverage.cov \
 -coverpkg ${COVER_PKG_LIST} \
 ${PACKAGE_LIST}

# Display the global code coverage
TOTAL=$(go tool cover -func="${COVERAGE_DIR}"/coverage.cov | grep "total:")
TOTAL=$(echo "$TOTAL" | grep -o "\d\+[\.]*\d*[%]")
echo "===================="
echo "$MODULE" "$TOTAL"
echo "===================="

COLOR="#DC143C"
PERCENT=$(echo "$TOTAL" | grep -o "\d\+" | head -1)

if [ "$PERCENT" -gt "49" ] ; then
  COLOR="#DFB317"
fi

if [ "$PERCENT" -gt "69" ] ; then
  COLOR="#00FF00"
fi

cat << EOF > "${COVERAGE_DIR}"/coverage.svg
<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" width="99" height="20">
    <linearGradient id="b" x2="0" y2="100%">
        <stop offset="0" stop-color="#bbb" stop-opacity=".1"/>
        <stop offset="1" stop-opacity=".1"/>
    </linearGradient>
    <mask id="a">
        <rect width="99" height="20" rx="3" fill="#fff"/>
    </mask>
    <g mask="url(#a)">
        <path fill="#555" d="M0 0h60v20H0z"/>
        <path fill="${COLOR}" d="M60 0h40v20H60z"/>
        <path fill="url(#b)" d="M0 0h99v20H0z"/>
    </g>
    <g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11">
        <text x="31.5" y="15" fill="#010101" fill-opacity=".3">coverage</text>
        <text x="31.5" y="14">Coverage</text>
        <text x="80" y="15" fill="#010101" fill-opacity=".3">${TOTAL}</text>
        <text x="80" y="14">${TOTAL}</text>
    </g>
</svg>
EOF

# If needed, generate HTML report
if [ "$1" == "html" ]; then
    go tool cover -html="${COVERAGE_DIR}"/coverage.cov -o "${COVERAGE_DIR}"/index.html
fi

BRANCH=`git rev-parse --abbrev-ref HEAD`;

if [[ "$1" == "html" ]] &&  [[ "$2" == "send" ]] &&  [[ "$BRANCH" == "dev" ]]; then
COVERAGE=$(echo "$TOTAL" | grep -o "\d\+[\.]*\d*" | head -1)
    curl --request POST $3'/api/v1/upload' \
    --form 'file=@"'$COVERAGE_DIR'/index.html"' \
    --form 'module="'$APPLICATION'"' \
    --form 'coverage="'$COVERAGE'"'
fi

