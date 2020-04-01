
# ch-server
Simple implemetation  of tusd server.
tus is a protocol based on HTTP for resumable file uploads. Resumable means that an upload can be interrupted at any moment and can be resumed without re-uploading the previous data again. An interruption may happen willingly, if the user wants to pause, or by accident in case of an network issue or server outage.

tusd is the official reference implementation of the tus resumable upload protocol. The protocol specifies a flexible method to upload files to remote servers using HTTP. The special feature is the ability to pause and resume uploads at any moment allowing to continue seamlessly after e.g. network interruptions.