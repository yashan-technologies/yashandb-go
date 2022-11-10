set(YACAPI_TEST
    ${PROJECT_SOURCE_DIR}/src/yapi_cli.c
    ${PROJECT_SOURCE_DIR}/src/yapi_conn.c
    ${PROJECT_SOURCE_DIR}/src/yapi_env.c
    ${PROJECT_SOURCE_DIR}/src/yapi_err.c
    ${PROJECT_SOURCE_DIR}/src/yapi_infr.c
    ${PROJECT_SOURCE_DIR}/src/yapi_stmt.c
    ${PROJECT_SOURCE_DIR}/src/yapi_var.c
    yapi_test.c
   )

add_executable(yasapitest ${YACAPI_TEST})
target_include_directories(yasapitest PUBLIC ${PROJECT_SOURCE_DIR}/src ${PROJECT_SOURCE_DIR}/include)