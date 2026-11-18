# Screen Observation Example

This example is read-only: it captures the current AutoGo-host screen, optionally saves a PNG to the host script directory, and prints recognized text with center coordinates. It performs no touch, app launch, network request, or form submission.

The checked-in AutoGo Go packages are local API stubs, so a local build only validates imports and signatures. Actual screen pixels/OCR require deployment and execution inside the AutoGo iOS host on an authorized device. Treat the generated PNG and OCR output as potentially sensitive; do not commit them.
