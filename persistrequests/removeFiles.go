package persistrequests

// RemoveRequest generates a filename from the given parameters the user has provided and will try to delete that request if it exists. An error will be returned if it doesn't exist
func (j JSONPersist) RemoveRequest(method, route string, r Remover) error {

	filename := createFilename(method, route)

	return r.Remove(j.RequestDirectory + filename)
}


