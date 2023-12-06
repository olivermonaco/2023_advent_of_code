#!/bin/bash
PROJECT_FOLDER=$(realpath $(dirname $(dirname "${BASH_SOURCE[0]}")))

YEAR="2023"
PACKAGE_REPLACE_STR="package_name"
DAY_REPLACE_STR="day_number"

EXAMPLE_FILES_FOLDER=$PROJECT_FOLDER/utils/example_files
MAIN_GO_TEMPLATE_FILENAME="main.go_"
PACKAGE_TEMPLATE_FILENAME="package_name.go_"
TEST_FILE_TEMPLATE_FILENAME="package_name_test.go_"
ADVENT_CHALLENGE_URL="https://adventofcode.com"


validate_args() {
  min_args=$1
  max_args=$2
  num_args=$3
  if ((num_args < min_args)); then
    echo not enough arguments. number must be at least $min_args and no more than $max_args
    return
  elif ((num_args > max_args)); then
    echo too many arguments. number must be at least $min_args and no more than $max_args
    return
  fi
  return
}

make_challenge() {
    invalid_num_args=$(validate_args 2 2 $#)
    if [ -n "$invalid_num_args" ]; then
        # invalid num args
        echo $invalid_num_args
        return
    fi
    # challenge name
    CHALLENGE_NAME=$1

    # calendar day
    DAY=$2

    # create date string with input day and fixed month and year
    DATE_STR="$YEAR-12-$DAY"

    # format date
    FORMATTED_DATE=$(date -j -f "%Y-%m-%d" "$DATE_STR" "+%d%m%y")

    TOP_FOLDER_NAME=$FORMATTED_DATE"_"$CHALLENGE_NAME
    TOP_FOLDER=$PROJECT_FOLDER/challenges/$TOP_FOLDER_NAME
    PACKAGE_FOLDER=$TOP_FOLDER/$CHALLENGE_NAME

    echo $TOP_FOLDER
    cd $PROJECT_FOLDER
    mkdir -p $PACKAGE_FOLDER/test_files

    # copy over main.go file, replace package_name with input
    cp $EXAMPLE_FILES_FOLDER/$MAIN_GO_TEMPLATE_FILENAME $TOP_FOLDER
    mv $TOP_FOLDER/$MAIN_GO_TEMPLATE_FILENAME $TOP_FOLDER/main.go


    # read file and replace string
    sed -i '' "s/$PACKAGE_REPLACE_STR/$CHALLENGE_NAME/g" $TOP_FOLDER/main.go
    sed -i '' "s/$PACKAGE_REPLACE_STR/$CHALLENGE_NAME/g" $TOP_FOLDER/main.go
    # include the url to the challenge
    ADVENT_DAY_URL_STR="// ${ADVENT_CHALLENGE_URL}/$YEAR/day/$DAY"
    sed -i '' "1i \\
    $ADVENT_DAY_URL_STR
        " $TOP_FOLDER/main.go

    # create folder for package
    mkdir -p $PACKAGE_FOLDER

    # copy over package .go file, replace package_name with input
    cp $EXAMPLE_FILES_FOLDER/$PACKAGE_TEMPLATE_FILENAME $PACKAGE_FOLDER
    mv $PACKAGE_FOLDER/$PACKAGE_TEMPLATE_FILENAME $PACKAGE_FOLDER/$CHALLENGE_NAME.go
    sed -i '' "s/$PACKAGE_REPLACE_STR/$CHALLENGE_NAME/g" "${PACKAGE_FOLDER}/${CHALLENGE_NAME}.go"


    # copy over package test file, replace package_name with input
    cp $EXAMPLE_FILES_FOLDER/$TEST_FILE_TEMPLATE_FILENAME $PACKAGE_FOLDER
    mv $PACKAGE_FOLDER/$TEST_FILE_TEMPLATE_FILENAME ${PACKAGE_FOLDER}/${CHALLENGE_NAME}_test.go
    sed -i '' "s/$PACKAGE_REPLACE_STR/$CHALLENGE_NAME/g" "${PACKAGE_FOLDER}/${CHALLENGE_NAME}_test.go"

    # placeholder for the puzzle input
    touch $TOP_FOLDER/puzzle_input.txt
}

make_challenge "$@"