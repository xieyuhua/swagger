<?php

require('vendor/autoload.php');
@$swagger = \Swagger\scan('Examples');
//var_dump(json_encode($swagger,JSON_UNESCAPED_UNICODE)) ;die;
header('Content-Type: application/json');
//echo $swagger;开辟
@file_put_contents('swagger.json',json_encode($swagger,JSON_UNESCAPED_UNICODE));
