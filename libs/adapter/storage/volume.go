package storage

import (
	"github.com/lastbackend/lastbackend/libs/interface/storage"
	"github.com/lastbackend/lastbackend/libs/model"
	r "gopkg.in/dancannon/gorethink.v2"
	"time"
)

const VolumeTable string = "volumes"

// Volume Service type for interface in interfaces folder
type VolumeStorage struct {
	Session *r.Session
	storage.IVolume
}

func (s *VolumeStorage) GetByID(user, id string) (*model.Volume, error) {

	var err error
	var volume = new(model.Volume)
	var volume_filter = map[string]interface{}{
		"id":   id,
		"user": user,
	}

	res, err := r.Table(VolumeTable).Filter(volume_filter).Run(s.Session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	if res.IsNil() {
		return nil, nil
	}

	res.One(volume)

	return volume, nil
}

func (s *VolumeStorage) ListByProject(id string) (*model.VolumeList, error) {

	var err error
	var volumes = new(model.VolumeList)
	var volume_filter = r.Row.Field("project").Eq(id)

	res, err := r.Table(VolumeTable).Filter(volume_filter).Run(s.Session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	if res.IsNil() {
		return nil, nil
	}

	res.All(volumes)

	return volumes, nil
}

// Insert new volume into storage
func (s *VolumeStorage) Insert(volume *model.Volume) (*model.Volume, error) {

	var err error
	var opts = r.InsertOpts{ReturnChanges: true}

	volume.Created = time.Now()
	volume.Updated = time.Now()

	res, err := r.Table(VolumeTable).Insert(volume, opts).RunWrite(s.Session)
	if err != nil {
		return nil, err
	}

	volume.ID = res.GeneratedKeys[0]

	return volume, nil
}

// Remove build model
func (s *VolumeStorage) Remove(id string) error {

	var err error
	var opts = r.DeleteOpts{ReturnChanges: true}

	_, err = r.Table(VolumeTable).Get(id).Delete(opts).RunWrite(s.Session)
	if err != nil {
		return err
	}

	return nil
}

func newVolumeStorage(session *r.Session) *VolumeStorage {
	r.TableCreate(VolumeTable, r.TableCreateOpts{}).Run(session)
	s := new(VolumeStorage)
	s.Session = session
	return s
}
