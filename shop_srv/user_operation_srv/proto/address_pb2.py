# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: address.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\raddress.proto\x1a\x1bgoogle/protobuf/empty.proto\"\x99\x01\n\x0e\x41\x64\x64ressRequest\x12\n\n\x02id\x18\x01 \x01(\x05\x12\x0e\n\x06userId\x18\x02 \x01(\x05\x12\x10\n\x08province\x18\x03 \x01(\t\x12\x0c\n\x04\x63ity\x18\x04 \x01(\t\x12\x10\n\x08\x64istrict\x18\x05 \x01(\t\x12\x0f\n\x07\x61\x64\x64ress\x18\x06 \x01(\t\x12\x12\n\nsignerName\x18\x07 \x01(\t\x12\x14\n\x0csignerMobile\x18\x08 \x01(\t\"\x9a\x01\n\x0f\x41\x64\x64ressResponse\x12\n\n\x02id\x18\x01 \x01(\x05\x12\x0e\n\x06userId\x18\x02 \x01(\x05\x12\x10\n\x08province\x18\x03 \x01(\t\x12\x0c\n\x04\x63ity\x18\x04 \x01(\t\x12\x10\n\x08\x64istrict\x18\x05 \x01(\t\x12\x0f\n\x07\x61\x64\x64ress\x18\x06 \x01(\t\x12\x12\n\nsignerName\x18\x07 \x01(\t\x12\x14\n\x0csignerMobile\x18\x08 \x01(\t\"D\n\x13\x41\x64\x64ressListResponse\x12\r\n\x05total\x18\x01 \x01(\x05\x12\x1e\n\x04\x64\x61ta\x18\x02 \x03(\x0b\x32\x10.AddressResponse2\xea\x01\n\x07\x41\x64\x64ress\x12\x37\n\x0eGetAddressList\x12\x0f.AddressRequest\x1a\x14.AddressListResponse\x12\x32\n\rCreateAddress\x12\x0f.AddressRequest\x1a\x10.AddressResponse\x12\x38\n\rDeleteAddress\x12\x0f.AddressRequest\x1a\x16.google.protobuf.Empty\x12\x38\n\rUpdateAddress\x12\x0f.AddressRequest\x1a\x16.google.protobuf.EmptyB\nZ\x08./;protob\x06proto3')



_ADDRESSREQUEST = DESCRIPTOR.message_types_by_name['AddressRequest']
_ADDRESSRESPONSE = DESCRIPTOR.message_types_by_name['AddressResponse']
_ADDRESSLISTRESPONSE = DESCRIPTOR.message_types_by_name['AddressListResponse']
AddressRequest = _reflection.GeneratedProtocolMessageType('AddressRequest', (_message.Message,), {
  'DESCRIPTOR' : _ADDRESSREQUEST,
  '__module__' : 'address_pb2'
  # @@protoc_insertion_point(class_scope:AddressRequest)
  })
_sym_db.RegisterMessage(AddressRequest)

AddressResponse = _reflection.GeneratedProtocolMessageType('AddressResponse', (_message.Message,), {
  'DESCRIPTOR' : _ADDRESSRESPONSE,
  '__module__' : 'address_pb2'
  # @@protoc_insertion_point(class_scope:AddressResponse)
  })
_sym_db.RegisterMessage(AddressResponse)

AddressListResponse = _reflection.GeneratedProtocolMessageType('AddressListResponse', (_message.Message,), {
  'DESCRIPTOR' : _ADDRESSLISTRESPONSE,
  '__module__' : 'address_pb2'
  # @@protoc_insertion_point(class_scope:AddressListResponse)
  })
_sym_db.RegisterMessage(AddressListResponse)

_ADDRESS = DESCRIPTOR.services_by_name['Address']
if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'Z\010./;proto'
  _ADDRESSREQUEST._serialized_start=47
  _ADDRESSREQUEST._serialized_end=200
  _ADDRESSRESPONSE._serialized_start=203
  _ADDRESSRESPONSE._serialized_end=357
  _ADDRESSLISTRESPONSE._serialized_start=359
  _ADDRESSLISTRESPONSE._serialized_end=427
  _ADDRESS._serialized_start=430
  _ADDRESS._serialized_end=664
# @@protoc_insertion_point(module_scope)
