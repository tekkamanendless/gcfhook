# gcfhook
Google Cloud Function hook for [logrus](https://github.com/sirupsen/logrus).

If you want to use "structured logging" (where you basically output JSON messages to standard output), then please see [gcfstructuredloghook](https://github.com/tekkamanendless/gcfstructuredloghook).

## Warning
This was originally intended to be used with Go 1.11 Cloud Functions, where the following environment variables were automatically defined:

* `FUNCTION_NAME`; this is the name of the Cloud Function.
* `FUNCTION_REGION`; this is the region of the Cloud Function.
* `GCP_PROJECT`; this is the Google Cloud project ID.

However, as of Go 1.13 Cloud Functions, these environment variables are no longer provided.
You may, of course, define them yourself and continue to use this package.

If you want to use "structured logging" (where you basically output JSON messages to standard output), then please see [gcfstructuredloghook](https://github.com/tekkamanendlessgcfstructuredloghook).

## Usage
This example shows how to set up logrus for Google Cloud Functions.
If the function is running locally (for development, etc.), then it will not use the hook.
Otherwise, it will set up the hook and "nullify" the normal `stderr` output using a `NullFormatter` (this prevents the log messges from appearing twice in the GCF logs).

```
// CloudFunction is an HTTP Cloud Function with a request parameter.
func CloudFunction(w http.ResponseWriter, r *http.Request) {
	log := logrus.New()

	if value := os.Getenv("GCP_PROJECT"); value == "" {
		log.Infof("GCP_PROJECT is not set; falling back to normal logging.")
	} else {
		hook, err := gcfhook.NewForRequest(r)
		if err != nil {
			log.Errorf("Could not set up gcfhook: %v", err)
		}

		if hook != nil {
			// Flush the logging entries when we're done.
			defer hook.Flush()

			// Add the hook.
			log.AddHook(hook)
			// Nullify the console output; we don't want to duplicate it.
			log.SetFormatter(&gcfhook.NullFormatter{})
		}
	}

	log.Infof("This is an info message.")
	log.Warnf("This is a warning message.")
	log.Errorf("This is an error message.")

	// YOUR CLOUD FUNCTION LOGIC HERE
}

```

