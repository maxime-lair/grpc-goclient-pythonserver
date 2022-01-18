# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: server.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='server.proto',
  package='server',
  syntax='proto3',
  serialized_options=None,
  create_key=_descriptor._internal_create_key,
  serialized_pb=b'\n\x0cserver.proto\x12\x06server\"\x1c\n\nSocketTree\x12\x0e\n\x06\x63hoice\x18\x01 \x01(\t\"+\n\x0cSocketFamily\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\x05\")\n\nSocketType\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\x05\"-\n\x0eSocketProtocol\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\x05\x32\xde\x01\n\x0bSocketGuide\x12\x43\n\x13GetSocketFamilyList\x12\x12.server.SocketTree\x1a\x14.server.SocketFamily\"\x00\x30\x01\x12\x41\n\x11GetSocketTypeList\x12\x14.server.SocketFamily\x1a\x12.server.SocketType\"\x00\x30\x01\x12G\n\x15GetSocketProtocolList\x12\x12.server.SocketType\x1a\x16.server.SocketProtocol\"\x00\x30\x01\x62\x06proto3'
)




_SOCKETTREE = _descriptor.Descriptor(
  name='SocketTree',
  full_name='server.SocketTree',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='choice', full_name='server.SocketTree.choice', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=24,
  serialized_end=52,
)


_SOCKETFAMILY = _descriptor.Descriptor(
  name='SocketFamily',
  full_name='server.SocketFamily',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='name', full_name='server.SocketFamily.name', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='value', full_name='server.SocketFamily.value', index=1,
      number=2, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=54,
  serialized_end=97,
)


_SOCKETTYPE = _descriptor.Descriptor(
  name='SocketType',
  full_name='server.SocketType',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='name', full_name='server.SocketType.name', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='value', full_name='server.SocketType.value', index=1,
      number=2, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=99,
  serialized_end=140,
)


_SOCKETPROTOCOL = _descriptor.Descriptor(
  name='SocketProtocol',
  full_name='server.SocketProtocol',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='name', full_name='server.SocketProtocol.name', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='value', full_name='server.SocketProtocol.value', index=1,
      number=2, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=142,
  serialized_end=187,
)

DESCRIPTOR.message_types_by_name['SocketTree'] = _SOCKETTREE
DESCRIPTOR.message_types_by_name['SocketFamily'] = _SOCKETFAMILY
DESCRIPTOR.message_types_by_name['SocketType'] = _SOCKETTYPE
DESCRIPTOR.message_types_by_name['SocketProtocol'] = _SOCKETPROTOCOL
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

SocketTree = _reflection.GeneratedProtocolMessageType('SocketTree', (_message.Message,), {
  'DESCRIPTOR' : _SOCKETTREE,
  '__module__' : 'server_pb2'
  # @@protoc_insertion_point(class_scope:server.SocketTree)
  })
_sym_db.RegisterMessage(SocketTree)

SocketFamily = _reflection.GeneratedProtocolMessageType('SocketFamily', (_message.Message,), {
  'DESCRIPTOR' : _SOCKETFAMILY,
  '__module__' : 'server_pb2'
  # @@protoc_insertion_point(class_scope:server.SocketFamily)
  })
_sym_db.RegisterMessage(SocketFamily)

SocketType = _reflection.GeneratedProtocolMessageType('SocketType', (_message.Message,), {
  'DESCRIPTOR' : _SOCKETTYPE,
  '__module__' : 'server_pb2'
  # @@protoc_insertion_point(class_scope:server.SocketType)
  })
_sym_db.RegisterMessage(SocketType)

SocketProtocol = _reflection.GeneratedProtocolMessageType('SocketProtocol', (_message.Message,), {
  'DESCRIPTOR' : _SOCKETPROTOCOL,
  '__module__' : 'server_pb2'
  # @@protoc_insertion_point(class_scope:server.SocketProtocol)
  })
_sym_db.RegisterMessage(SocketProtocol)



_SOCKETGUIDE = _descriptor.ServiceDescriptor(
  name='SocketGuide',
  full_name='server.SocketGuide',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  create_key=_descriptor._internal_create_key,
  serialized_start=190,
  serialized_end=412,
  methods=[
  _descriptor.MethodDescriptor(
    name='GetSocketFamilyList',
    full_name='server.SocketGuide.GetSocketFamilyList',
    index=0,
    containing_service=None,
    input_type=_SOCKETTREE,
    output_type=_SOCKETFAMILY,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='GetSocketTypeList',
    full_name='server.SocketGuide.GetSocketTypeList',
    index=1,
    containing_service=None,
    input_type=_SOCKETFAMILY,
    output_type=_SOCKETTYPE,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='GetSocketProtocolList',
    full_name='server.SocketGuide.GetSocketProtocolList',
    index=2,
    containing_service=None,
    input_type=_SOCKETTYPE,
    output_type=_SOCKETPROTOCOL,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
])
_sym_db.RegisterServiceDescriptor(_SOCKETGUIDE)

DESCRIPTOR.services_by_name['SocketGuide'] = _SOCKETGUIDE

# @@protoc_insertion_point(module_scope)