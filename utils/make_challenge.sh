#!/bin/bash
snake_to_camel() {
    echo $1 | awk 'BEGIN{FS="_";OFS=""} {for(i=1;i<=NF;i++){$i=toupper(substr($i,1,1)) substr($i,2)}} 1'
}

PROJECT_FOLDER=$(realpath $(dirname $(dirname "${BASH_SOURCE[0]}")))

YEAR="2023"
PACKAGE_SNAKE_CASE_REPLACE_STR="package_name"
PACKAGE_CAMEL_CASE_REPLACE_STR=$(snake_to_camel $PACKAGE_SNAKE_CASE_REPLACE_STR)

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
    CHALLENGE_SNAKE_CASE=$1
    CHALLENGE_CAMEL_CASE=$(snake_to_camel $CHALLENGE_SNAKE_CASE)

    # calendar day
    DAY=$2

    # create date string with input day and fixed month and year
    DATE_STR="$YEAR-12-$DAY"

    # format date
    FORMATTED_DATE=$(date -j -f "%Y-%m-%d" "$DATE_STR" "+%d%m%y")

    TOP_FOLDER_NAME=$FORMATTED_DATE"_"$CHALLENGE_SNAKE_CASE
    TOP_FOLDER=$PROJECT_FOLDER/challenges/$TOP_FOLDER_NAME
    PACKAGE_FOLDER=$TOP_FOLDER/$CHALLENGE_SNAKE_CASE

    echo $TOP_FOLDER
    cd $PROJECT_FOLDER
    mkdir -p $PACKAGE_FOLDER/test_files

    # copy over main.go file, replace package_name with input
    cp $EXAMPLE_FILES_FOLDER/$MAIN_GO_TEMPLATE_FILENAME $TOP_FOLDER
    mv $TOP_FOLDER/$MAIN_GO_TEMPLATE_FILENAME $TOP_FOLDER/main.go


    # read file and replace string
    sed -i '' "s/$PACKAGE_SNAKE_CASE_REPLACE_STR/$CHALLENGE_SNAKE_CASE/g" $TOP_FOLDER/main.go
    # include the url to the challenge
    ADVENT_DAY_URL_STR="// ${ADVENT_CHALLENGE_URL}/$YEAR/day/$DAY"
    sed -i '' "1i \\
    $ADVENT_DAY_URL_STR
        " $TOP_FOLDER/main.go

    # create folder for package
    mkdir -p $PACKAGE_FOLDER

    # copy over package .go file, replace package_name with input
    cp $EXAMPLE_FILES_FOLDER/$PACKAGE_TEMPLATE_FILENAME $PACKAGE_FOLDER
    mv $PACKAGE_FOLDER/$PACKAGE_TEMPLATE_FILENAME $PACKAGE_FOLDER/$CHALLENGE_SNAKE_CASE.go
    sed -i '' "s/$PACKAGE_SNAKE_CASE_REPLACE_STR/$CHALLENGE_SNAKE_CASE/g" "${PACKAGE_FOLDER}/${CHALLENGE_SNAKE_CASE}.go"


    # copy over package test file, replace package_name with input
    cp $EXAMPLE_FILES_FOLDER/$TEST_FILE_TEMPLATE_FILENAME $PACKAGE_FOLDER
    mv $PACKAGE_FOLDER/$TEST_FILE_TEMPLATE_FILENAME ${PACKAGE_FOLDER}/${CHALLENGE_SNAKE_CASE}_test.go
    sed -i '' "s/$PACKAGE_SNAKE_CASE_REPLACE_STR/$CHALLENGE_SNAKE_CASE/g" "${PACKAGE_FOLDER}/${CHALLENGE_SNAKE_CASE}_test.go"
    sed -i '' "s/$PACKAGE_CAMEL_CASE_REPLACE_STR/$CHALLENGE_CAMEL_CASE/g" "${PACKAGE_FOLDER}/${CHALLENGE_SNAKE_CASE}_test.go"

    # placeholder for the puzzle input
    touch $PACKAGE_FOLDER/puzzle_input.txt
}

make_challenge "$@"