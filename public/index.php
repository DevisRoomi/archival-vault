<?php
/**
 * Archival Vault Herd Gateway
 * Proxies /api requests to Go Backend (port 8085)
 * Serves SPA index.html for frontend routes
 */

$uri = $_SERVER['REQUEST_URI'];
$backendUrl = 'http://127.0.0.1:8085';

// If static file in public exists (and is not index.php), let webserver serve it
$publicFile = __DIR__ . parse_url($uri, PHP_URL_PATH);
if ($uri !== '/' && file_exists($publicFile) && !is_dir($publicFile) && basename($publicFile) !== 'index.php') {
    return false; // let Nginx serve directly
}

// Proxy API requests to Golang backend
if (str_starts_with(parse_url($uri, PHP_URL_PATH), '/api')) {
    $target = $backendUrl . $uri;
    $ch = curl_init($target);
    
    // Forward method
    curl_setopt($ch, CURLOPT_CUSTOMREQUEST, $_SERVER['REQUEST_METHOD']);
    
    // Forward request body
    $body = file_get_contents('php://input');
    if (!empty($body)) {
        curl_setopt($ch, CURLOPT_POSTFIELDS, $body);
    }
    
    // Forward request headers
    $headers = [];
    foreach (getallheaders() as $name => $value) {
        if (!in_array(strtolower($name), ['host', 'content-length'])) {
            $headers[] = "$name: $value";
        }
    }
    curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch, CURLOPT_HEADER, true);
    curl_setopt($ch, CURLOPT_TIMEOUT, 15);
    
    $response = curl_exec($ch);
    
    if (curl_errno($ch)) {
        http_response_code(502);
        header('Content-Type: application/json');
        echo json_encode([
            'success' => false,
            'error' => 'Gateway error: Unable to reach Go backend at ' . $backendUrl . '. Make sure backend is running.'
        ]);
        curl_close($ch);
        exit;
    }
    
    $headerSize = curl_getinfo($ch, CURLINFO_HEADER_SIZE);
    $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
    $resHeaders = substr($response, 0, $headerSize);
    $resBody = substr($response, $headerSize);
    
    http_response_code($httpCode);
    
    // Forward response headers
    foreach (explode("\r\n", $resHeaders) as $header) {
        if (!empty($header) && !str_starts_with(strtolower($header), 'transfer-encoding:')) {
            header($header, false);
        }
    }
    
    echo $resBody;
    curl_close($ch);
    exit;
}

// Otherwise serve SPA index.html
readfile(__DIR__ . '/index.html');
