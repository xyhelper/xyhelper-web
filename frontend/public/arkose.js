var arkoseReady = false; // Flag for checking if CAPI has called onReady
var arkoseComplete = false; // Flag to signify the challenge being completed
var arkoseResetting = false; // Flag to signify when the challenge is being reset
var arkose; // Arkose Enforcement Object
var arkoseRetry = 0; // Counter for Retries

// The following variables can be changed
var arkoseMaxRetryCount = 2; // The number of retries to perform when there is an error 
var arkoseCookieName = 'arkoseToken'; // The name of the cookie that the Arkose token will be stored in
var arkoseErrorCookieName = 'arkoseError'; // The name of the cookie that an Arkose error will be stored in
var arkoseCookieLife = 5 * 60 * 1000;  // The length of time that the cookie should be active for, 5 mins is the default
var buttonSelector = '#submitButton'; // The querySelector string used for selecting the required button to protect

var button = document.querySelector(buttonSelector);
setupButton();

/**
* Sets up the events for hijacking the button click
*/
function setupButton() {
    if (button) {
        arkoseComplete = false;
        if (!arkoseReady) {
            button.setAttribute('disabled', true);
        }
        button.addEventListener('click', function (event) {
            if (!arkoseReady) {
                event.preventDefault();
                return;
            }
            if (!arkoseComplete) {
                event.preventDefault();
                arkose.run();
                return;
            }
        })
    }
}

/**
* Checks the current status of the Arkose Labs platform and invokes a callback with the result
* @param  {function} callback The function to call once the status check has been performed
*/
function checkArkoseStatus(callback) {
    try {
        var xhr = new XMLHttpRequest();
        xhr.open('GET', 'https://status.arkoselabs.com/api/v2/status.json', false);
        xhr.onreadystatechange = function() {
            if (xhr.readyState == XMLHttpRequest.DONE) {
                if (this.status == 200) {
                    var res = JSON.parse(xhr.responseText);
                    var status = res.status.indicator;
                    callback(!(status === 'critical'));
                    return;
                }
                callback(false);
            }
        };
        xhr.send(null);
    } catch (error) {
        callback(false);
    }
}

/**
* In the case of an error reset the arkose cookie and setup a new cookie that tracks the error
* @param  {string} error The error string to set in the error cookie
*/
function handleError(error) {
    arkoseComplete = true;
    document.cookie = arkoseCookieName + '=;expires=' + new Date(Date.now() + arkoseCookieLife).toUTCString() + '; path=/;';
    document.cookie = arkoseErrorCookieName + '=' + error + ';expires=' + new Date(Date.now() + arkoseCookieLife).toUTCString() + '; path=/;';                               
}  

 /**
* The callback thats called after the Arkose Labs script has been loaded
* @param  {object} myEnforcement The Arkose Labs enforcement object
*/
function setupEnforcement(myEnforcement) {
    arkose = myEnforcement;
    arkose.setConfig({
        onReady: function() {
            arkoseReady = true;
            if (button) {
                button.removeAttribute('disabled');
            }
            if (arkoseResetting) {
                arkoseResetting = false;
                arkose.run();
            }
            document.cookie = arkoseCookieName + '=' + '=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
            document.cookie = arkoseErrorCookieName + '=' + '=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
        },
        onCompleted: function(response) {
            arkoseComplete = true;
            if (response.token){
                document.cookie = arkoseCookieName + '=' + response.token + ';expires=' + new Date(Date.now() + arkoseCookieLife).toUTCString() + '; path=/;';
            } else {
                handleError('TOKEN_MISSING');
            }
            button.click();
        },
        onError: function(response) {
            checkArkoseStatus(function (isHealthy) {    
                if (isHealthy && arkoseRetry < arkoseMaxRetryCount) {
                    arkoseReady = false;
                    arkoseResetting = true;
                    arkose.reset();
                    arkoseRetry = arkoseRetry + 1;
                    return;
                }
                handleError(response.error ? response.error.error : 'error');
                button.removeAttribute('disabled');
                button.click();
            });
        }
    });
}