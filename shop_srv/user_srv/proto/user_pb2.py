# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: user.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\nuser.proto\x1a\x1bgoogle/protobuf/empty.proto\"A\n\x11PasswordCheckInfo\x12\x10\n\x08password\x18\x01 \x01(\t\x12\x1a\n\x12\x65ncrypted_password\x18\x02 \x01(\t\" \n\rCheckResponse\x12\x0f\n\x07success\x18\x01 \x01(\x08\"&\n\x08PageInfo\x12\n\n\x02pn\x18\x01 \x01(\r\x12\x0e\n\x06p_size\x18\x02 \x01(\r\"B\n\x10UserListResponse\x12\r\n\x05total\x18\x01 \x01(\x05\x12\x1f\n\x04\x64\x61ta\x18\x02 \x03(\x0b\x32\x11.UserInfoResponse\"\x82\x01\n\x10UserInfoResponse\x12\n\n\x02id\x18\x01 \x01(\x05\x12\x10\n\x08password\x18\x02 \x01(\t\x12\x0e\n\x06mobile\x18\x03 \x01(\t\x12\x10\n\x08nickname\x18\x04 \x01(\t\x12\x10\n\x08\x62irthday\x18\x05 \x01(\x04\x12\x0e\n\x06gender\x18\x06 \x01(\t\x12\x0c\n\x04role\x18\x07 \x01(\x05\"Q\n\x0eUpdateUserInfo\x12\n\n\x02id\x18\x01 \x01(\x05\x12\x11\n\tnick_name\x18\x02 \x01(\t\x12\x0e\n\x06gender\x18\x03 \x01(\t\x12\x10\n\x08\x62irthday\x18\x04 \x01(\x04\"\x17\n\tIdRequest\x12\n\n\x02id\x18\x01 \x01(\x05\"E\n\x0e\x43reateUserInfo\x12\x11\n\tnick_name\x18\x01 \x01(\t\x12\x10\n\x08password\x18\x02 \x01(\t\x12\x0e\n\x06mobile\x18\x03 \x01(\t\"\x7f\n\rMobileRequest\x12\n\n\x02id\x18\x01 \x01(\t\x12\x10\n\x08password\x18\x02 \x01(\t\x12\x0e\n\x06mobile\x18\x03 \x01(\t\x12\x10\n\x08nickname\x18\x04 \x01(\t\x12\x10\n\x08\x62irthday\x18\x05 \x01(\t\x12\x0e\n\x06gender\x18\x06 \x01(\t\x12\x0c\n\x04role\x18\x07 \x01(\x05\x32\xb5\x02\n\x04User\x12+\n\x0bGetUserList\x12\t.PageInfo\x1a\x11.UserListResponse\x12\x34\n\x0fGetUserByMobile\x12\x0e.MobileRequest\x1a\x11.UserInfoResponse\x12,\n\x0bGetUserById\x12\n.IdRequest\x1a\x11.UserInfoResponse\x12\x30\n\nCreateUser\x12\x0f.CreateUserInfo\x1a\x11.UserInfoResponse\x12\x35\n\nUpdateUser\x12\x0f.UpdateUserInfo\x1a\x16.google.protobuf.Empty\x12\x33\n\rCheckPassWord\x12\x12.PasswordCheckInfo\x1a\x0e.CheckResponseB\nZ\x08./;protob\x06proto3')



_PASSWORDCHECKINFO = DESCRIPTOR.message_types_by_name['PasswordCheckInfo']
_CHECKRESPONSE = DESCRIPTOR.message_types_by_name['CheckResponse']
_PAGEINFO = DESCRIPTOR.message_types_by_name['PageInfo']
_USERLISTRESPONSE = DESCRIPTOR.message_types_by_name['UserListResponse']
_USERINFORESPONSE = DESCRIPTOR.message_types_by_name['UserInfoResponse']
_UPDATEUSERINFO = DESCRIPTOR.message_types_by_name['UpdateUserInfo']
_IDREQUEST = DESCRIPTOR.message_types_by_name['IdRequest']
_CREATEUSERINFO = DESCRIPTOR.message_types_by_name['CreateUserInfo']
_MOBILEREQUEST = DESCRIPTOR.message_types_by_name['MobileRequest']
PasswordCheckInfo = _reflection.GeneratedProtocolMessageType('PasswordCheckInfo', (_message.Message,), {
  'DESCRIPTOR' : _PASSWORDCHECKINFO,
  '__module__' : 'user_pb2'
  # @@protoc_insertion_point(class_scope:PasswordCheckInfo)
  })
_sym_db.RegisterMessage(PasswordCheckInfo)

CheckResponse = _reflection.GeneratedProtocolMessageType('CheckResponse', (_message.Message,), {
  'DESCRIPTOR' : _CHECKRESPONSE,
  '__module__' : 'user_pb2'
  # @@protoc_insertion_point(class_scope:CheckResponse)
  })
_sym_db.RegisterMessage(CheckResponse)

PageInfo = _reflection.GeneratedProtocolMessageType('PageInfo', (_message.Message,), {
  'DESCRIPTOR' : _PAGEINFO,
  '__module__' : 'user_pb2'
  # @@protoc_insertion_point(class_scope:PageInfo)
  })
_sym_db.RegisterMessage(PageInfo)

UserListResponse = _reflection.GeneratedProtocolMessageType('UserListResponse', (_message.Message,), {
  'DESCRIPTOR' : _USERLISTRESPONSE,
  '__module__' : 'user_pb2'
  # @@protoc_insertion_point(class_scope:UserListResponse)
  })
_sym_db.RegisterMessage(UserListResponse)

UserInfoResponse = _reflection.GeneratedProtocolMessageType('UserInfoResponse', (_message.Message,), {
  'DESCRIPTOR' : _USERINFORESPONSE,
  '__module__' : 'user_pb2'
  # @@protoc_insertion_point(class_scope:UserInfoResponse)
  })
_sym_db.RegisterMessage(UserInfoResponse)

UpdateUserInfo = _reflection.GeneratedProtocolMessageType('UpdateUserInfo', (_message.Message,), {
  'DESCRIPTOR' : _UPDATEUSERINFO,
  '__module__' : 'user_pb2'
  # @@protoc_insertion_point(class_scope:UpdateUserInfo)
  })
_sym_db.RegisterMessage(UpdateUserInfo)

IdRequest = _reflection.GeneratedProtocolMessageType('IdRequest', (_message.Message,), {
  'DESCRIPTOR' : _IDREQUEST,
  '__module__' : 'user_pb2'
  # @@protoc_insertion_point(class_scope:IdRequest)
  })
_sym_db.RegisterMessage(IdRequest)

CreateUserInfo = _reflection.GeneratedProtocolMessageType('CreateUserInfo', (_message.Message,), {
  'DESCRIPTOR' : _CREATEUSERINFO,
  '__module__' : 'user_pb2'
  # @@protoc_insertion_point(class_scope:CreateUserInfo)
  })
_sym_db.RegisterMessage(CreateUserInfo)

MobileRequest = _reflection.GeneratedProtocolMessageType('MobileRequest', (_message.Message,), {
  'DESCRIPTOR' : _MOBILEREQUEST,
  '__module__' : 'user_pb2'
  # @@protoc_insertion_point(class_scope:MobileRequest)
  })
_sym_db.RegisterMessage(MobileRequest)

_USER = DESCRIPTOR.services_by_name['User']
if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'Z\010./;proto'
  _PASSWORDCHECKINFO._serialized_start=43
  _PASSWORDCHECKINFO._serialized_end=108
  _CHECKRESPONSE._serialized_start=110
  _CHECKRESPONSE._serialized_end=142
  _PAGEINFO._serialized_start=144
  _PAGEINFO._serialized_end=182
  _USERLISTRESPONSE._serialized_start=184
  _USERLISTRESPONSE._serialized_end=250
  _USERINFORESPONSE._serialized_start=253
  _USERINFORESPONSE._serialized_end=383
  _UPDATEUSERINFO._serialized_start=385
  _UPDATEUSERINFO._serialized_end=466
  _IDREQUEST._serialized_start=468
  _IDREQUEST._serialized_end=491
  _CREATEUSERINFO._serialized_start=493
  _CREATEUSERINFO._serialized_end=562
  _MOBILEREQUEST._serialized_start=564
  _MOBILEREQUEST._serialized_end=691
  _USER._serialized_start=694
  _USER._serialized_end=1003
# @@protoc_insertion_point(module_scope)
