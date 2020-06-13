package noodle

//GLEnum represents all available WebGL constants, prefixed with Gl and turned into UpperCamelCase. For example, DEPTH_BUFFER_BIT is now GlDepthBufferBit
type GLEnum = int
const (
//GlDepthBufferBit passed to <code>clear</code> to clear the current depth buffer.
GlDepthBufferBit GLEnum = 0x00000100
//GlStencilBufferBit passed to <code>clear</code> to clear the current stencil buffer.
GlStencilBufferBit = 0x00000400
//GlColorBufferBit passed to <code>clear</code> to clear the current color buffer.
GlColorBufferBit = 0x00004000
//GlPoints passed to <code>drawelements</code> or <code>drawarrays</code> to draw single points.
GlPoints = 0x0000
//GlLines passed to <code>drawelements</code> or <code>drawarrays</code> to draw lines. each vertex connects to the one after it.
GlLines = 0x0001
//GlLineLoop passed to <code>drawelements</code> or <code>drawarrays</code> to draw lines. each set of two vertices is treated as a separate line segment.
GlLineLoop = 0x0002
//GlLineStrip passed to <code>drawelements</code> or <code>drawarrays</code> to draw a connected group of line segments from the first vertex to the last.
GlLineStrip = 0x0003
//GlTriangles passed to <code>drawelements</code> or <code>drawarrays</code> to draw triangles. each set of three vertices creates a separate triangle.
GlTriangles = 0x0004
//GlTriangleStrip passed to <code>drawelements</code> or <code>drawarrays</code> to draw a connected group of triangles.
GlTriangleStrip = 0x0005
//GlTriangleFan passed to <code>drawelements</code> or <code>drawarrays</code> to draw a connected group of triangles. each vertex connects to the previous and the first vertex in the fan.
GlTriangleFan = 0x0006
//GlZero passed to <code>blendfunc</code> or <code>blendfuncseparate</code> to turn off a component.
GlZero = 0
//GlOne passed to <code>blendfunc</code> or <code>blendfuncseparate</code> to turn on a component.
GlOne = 1
//GlSrcColor passed to <code>blendfunc</code> or <code>blendfuncseparate</code> to multiply a component by the source elements color.
GlSrcColor = 0x0300
//GlOneMinusSrcColor passed to <code>blendfunc</code> or <code>blendfuncseparate</code> to multiply a component by one minus the source elements color.
GlOneMinusSrcColor = 0x0301
//GlSrcAlpha passed to <code>blendfunc</code> or <code>blendfuncseparate</code> to multiply a component by the source's alpha.
GlSrcAlpha = 0x0302
//GlOneMinusSrcAlpha passed to <code>blendfunc</code> or <code>blendfuncseparate</code> to multiply a component by one minus the source's alpha.
GlOneMinusSrcAlpha = 0x0303
//GlDstAlpha passed to <code>blendfunc</code> or <code>blendfuncseparate</code> to multiply a component by the destination's alpha.
GlDstAlpha = 0x0304
//GlOneMinusDstAlpha passed to <code>blendfunc</code> or <code>blendfuncseparate</code> to multiply a component by one minus the destination's alpha.
GlOneMinusDstAlpha = 0x0305
//GlDstColor passed to <code>blendfunc</code> or <code>blendfuncseparate</code> to multiply a component by the destination's color.
GlDstColor = 0x0306
//GlOneMinusDstColor passed to <code>blendfunc</code> or <code>blendfuncseparate</code> to multiply a component by one minus the destination's color.
GlOneMinusDstColor = 0x0307
//GlSrcAlphaSaturate passed to <code>blendfunc</code> or <code>blendfuncseparate</code> to multiply a component by the minimum of source's alpha or one minus the destination's alpha.
GlSrcAlphaSaturate = 0x0308
//GlConstantColor passed to <code>blendfunc</code> or <code>blendfuncseparate</code> to specify a constant color blend function.
GlConstantColor = 0x8001
//GlOneMinusConstantColor passed to <code>blendfunc</code> or <code>blendfuncseparate</code> to specify one minus a constant color blend function.
GlOneMinusConstantColor = 0x8002
//GlConstantAlpha passed to <code>blendfunc</code> or <code>blendfuncseparate</code> to specify a constant alpha blend function.
GlConstantAlpha = 0x8003
//GlOneMinusConstantAlpha passed to <code>blendfunc</code> or <code>blendfuncseparate</code> to specify one minus a constant alpha blend function.
GlOneMinusConstantAlpha = 0x8004
//GlFuncAdd passed to <code>blendequation</code> or <code>blendequationseparate</code> to set an addition blend function.
GlFuncAdd = 0x8006
//GlFuncSubtract passed to <code>blendequation</code> or <code>blendequationseparate</code> to specify a subtraction blend function (source - destination).
GlFuncSubtract = 0x800a
//GlFuncReverseSubtract passed to <code>blendequation</code> or <code>blendequationseparate</code> to specify a reverse subtraction blend function (destination - source).
GlFuncReverseSubtract = 0x800b
//GlBlendEquation passed to <code>getparameter</code> to get the current rgb blend function.
GlBlendEquation = 0x8009
//GlBlendEquationRgb passed to <code>getparameter</code> to get the current rgb blend function. same as blendEquation
GlBlendEquationRgb = 0x8009
//GlBlendEquationAlpha passed to <code>getparameter</code> to get the current alpha blend function. same as blendEquation
GlBlendEquationAlpha = 0x883d
//GlBlendDstRgb passed to <code>getparameter</code> to get the current destination rgb blend function.
GlBlendDstRgb = 0x80c8
//GlBlendSrcRgb passed to <code>getparameter</code> to get the current destination rgb blend function.
GlBlendSrcRgb = 0x80c9
//GlBlendDstAlpha passed to <code>getparameter</code> to get the current destination alpha blend function.
GlBlendDstAlpha = 0x80ca
//GlBlendSrcAlpha passed to <code>getparameter</code> to get the current source alpha blend function.
GlBlendSrcAlpha = 0x80cb
//GlBlendColor passed to <code>getparameter</code> to return a the current blend color.
GlBlendColor = 0x8005
//GlArrayBufferBinding passed to <code>getparameter</code> to get the array buffer binding.
GlArrayBufferBinding = 0x8894
//GlElementArrayBufferBinding passed to <code>getparameter</code> to get the current element array buffer.
GlElementArrayBufferBinding = 0x8895
//GlLineWidth passed to <code>getparameter</code> to get the current <code>linewidth</code> (set by the <code>linewidth</code> method).
GlLineWidth = 0x0b21
//GlAliasedPointSizeRange passed to <code>getparameter</code> to get the current size of a point drawn with <code>gl.points</code>
GlAliasedPointSizeRange = 0x846d
//GlAliasedLineWidthRange passed to <code>getparameter</code> to get the range of available widths for a line. returns a length-2 array with the lo value at 0, and hight at 1.
GlAliasedLineWidthRange = 0x846e
//GlCullFaceMode passed to <code>getparameter</code> to get the current value of <code>cullface</code>. should return <code>front</code>, <code>back</code>, or <code>frontAndBack</code>
GlCullFaceMode = 0x0b45
//GlFrontFace passed to <code>getparameter</code> to determine the current value of <code>frontface</code>. should return <code>cw</code> or <code>ccw</code>.
GlFrontFace = 0x0b46
//GlDepthRange passed to <code>getparameter</code> to return a length-2 array of floats giving the current depth range.
GlDepthRange = 0x0b70
//GlDepthWritemask passed to <code>getparameter</code> to determine if the depth write mask is enabled.
GlDepthWritemask = 0x0b72
//GlDepthClearValue passed to <code>getparameter</code> to determine the current depth clear value.
GlDepthClearValue = 0x0b73
//GlDepthFunc passed to <code>getparameter</code> to get the current depth function. returns <code>never</code>, <code>always</code>, <code>less</code>, <code>equal</code>, <code>lequal</code>, <code>greater</code>, <code>gequal</code>, or <code>notequal</code>.
GlDepthFunc = 0x0b74
//GlStencilClearValue passed to <code>getparameter</code> to get the value the stencil will be cleared to.
GlStencilClearValue = 0x0b91
//GlStencilFunc passed to <code>getparameter</code> to get the current stencil function. returns <code>never</code>, <code>always</code>, <code>less</code>, <code>equal</code>, <code>lequal</code>, <code>greater</code>, <code>gequal</code>, or <code>notequal</code>.
GlStencilFunc = 0x0b92
//GlStencilFail passed to <code>getparameter</code> to get the current stencil fail function. should return <code>keep</code>, <code>replace</code>, <code>incr</code>, <code>decr</code>, <code>invert</code>, <code>incrWrap</code>, or <code>decrWrap</code>.
GlStencilFail = 0x0b94
//GlStencilPassDepthFail passed to <code>getparameter</code> to get the current stencil fail function should the depth buffer test fail. should return <code>keep</code>, <code>replace</code>, <code>incr</code>, <code>decr</code>, <code>invert</code>, <code>incrWrap</code>, or <code>decrWrap</code>.
GlStencilPassDepthFail = 0x0b95
//GlStencilPassDepthPass passed to <code>getparameter</code> to get the current stencil fail function should the depth buffer test pass. should return keep, replace, incr, decr, invert, incrWrap, or decrWrap.
GlStencilPassDepthPass = 0x0b96
//GlStencilRef passed to <code>getparameter</code> to get the reference value used for stencil tests.
GlStencilRef = 0x0b97
//GlStencilValueMask &nbsp;
GlStencilValueMask = 0x0b93
//GlStencilWritemask &nbsp;
GlStencilWritemask = 0x0b98
//GlStencilBackFunc &nbsp;
GlStencilBackFunc = 0x8800
//GlStencilBackFail &nbsp;
GlStencilBackFail = 0x8801
//GlStencilBackPassDepthFail &nbsp;
GlStencilBackPassDepthFail = 0x8802
//GlStencilBackPassDepthPass &nbsp;
GlStencilBackPassDepthPass = 0x8803
//GlStencilBackRef &nbsp;
GlStencilBackRef = 0x8ca3
//GlStencilBackValueMask &nbsp;
GlStencilBackValueMask = 0x8ca4
//GlStencilBackWritemask &nbsp;
GlStencilBackWritemask = 0x8ca5
//GlViewport returns an <a href="/en-us/docs/web/javascript/reference/globalObjects/int32array" title="the int32array typed array represents an array of twos-complement 32-bit signed integers in the platform byte order. if control over byte order is needed, use dataview instead. the contents are initialized to 0. once established, you can reference elements in the array using the object's methods, or using standard array index syntax (that is, using bracket notation)."><code>int32array</code></a> with four elements for the current viewport dimensions.
GlViewport = 0x0ba2
//GlScissorBox returns an <a href="/en-us/docs/web/javascript/reference/globalObjects/int32array" title="the int32array typed array represents an array of twos-complement 32-bit signed integers in the platform byte order. if control over byte order is needed, use dataview instead. the contents are initialized to 0. once established, you can reference elements in the array using the object's methods, or using standard array index syntax (that is, using bracket notation)."><code>int32array</code></a> with four elements for the current scissor box dimensions.
GlScissorBox = 0x0c10
//GlColorClearValue &nbsp;
GlColorClearValue = 0x0c22
//GlColorWritemask &nbsp;
GlColorWritemask = 0x0c23
//GlUnpackAlignment &nbsp;
GlUnpackAlignment = 0x0cf5
//GlPackAlignment &nbsp;
GlPackAlignment = 0x0d05
//GlMaxTextureSize &nbsp;
GlMaxTextureSize = 0x0d33
//GlMaxViewportDims &nbsp;
GlMaxViewportDims = 0x0d3a
//GlSubpixelBits &nbsp;
GlSubpixelBits = 0x0d50
//GlRedBits &nbsp;
GlRedBits = 0x0d52
//GlGreenBits &nbsp;
GlGreenBits = 0x0d53
//GlBlueBits &nbsp;
GlBlueBits = 0x0d54
//GlAlphaBits &nbsp;
GlAlphaBits = 0x0d55
//GlDepthBits &nbsp;
GlDepthBits = 0x0d56
//GlStencilBits &nbsp;
GlStencilBits = 0x0d57
//GlPolygonOffsetUnits &nbsp;
GlPolygonOffsetUnits = 0x2a00
//GlPolygonOffsetFactor &nbsp;
GlPolygonOffsetFactor = 0x8038
//GlTextureBinding2d &nbsp;
GlTextureBinding2d = 0x8069
//GlSampleBuffers &nbsp;
GlSampleBuffers = 0x80a8
//GlSamples &nbsp;
GlSamples = 0x80a9
//GlSampleCoverageValue &nbsp;
GlSampleCoverageValue = 0x80aa
//GlSampleCoverageInvert &nbsp;
GlSampleCoverageInvert = 0x80ab
//GlCompressedTextureFormats &nbsp;
GlCompressedTextureFormats = 0x86a3
//GlVendor &nbsp;
GlVendor = 0x1f00
//GlRenderer &nbsp;
GlRenderer = 0x1f01
//GlVersion &nbsp;
GlVersion = 0x1f02
//GlImplementationColorReadType &nbsp;
GlImplementationColorReadType = 0x8b9a
//GlImplementationColorReadFormat &nbsp;
GlImplementationColorReadFormat = 0x8b9b
//GlBrowserDefaultWebgl &nbsp;
GlBrowserDefaultWebgl = 0x9244
//GlStaticDraw passed to <code>bufferdata</code> as a hint about whether the contents of the buffer are likely to be used often and not change often.
GlStaticDraw = 0x88e4
//GlStreamDraw passed to <code>bufferdata</code> as a hint about whether the contents of the buffer are likely to not be used often.
GlStreamDraw = 0x88e0
//GlDynamicDraw passed to <code>bufferdata</code> as a hint about whether the contents of the buffer are likely to be used often and change often.
GlDynamicDraw = 0x88e8
//GlArrayBuffer passed to <code>bindbuffer</code> or <code>bufferdata</code> to specify the type of buffer being used.
GlArrayBuffer = 0x8892
//GlElementArrayBuffer passed to <code>bindbuffer</code> or <code>bufferdata</code> to specify the type of buffer being used.
GlElementArrayBuffer = 0x8893
//GlBufferSize passed to <code>getbufferparameter</code> to get a buffer's size.
GlBufferSize = 0x8764
//GlBufferUsage passed to&nbsp;<code>getbufferparameter</code> to get the hint for the buffer passed in when it was created.
GlBufferUsage = 0x8765
//GlCurrentVertexAttrib passed to <code>getvertexattrib</code> to read back the current vertex attribute.
GlCurrentVertexAttrib = 0x8626
//GlVertexAttribArrayEnabled &nbsp;
GlVertexAttribArrayEnabled = 0x8622
//GlVertexAttribArraySize &nbsp;
GlVertexAttribArraySize = 0x8623
//GlVertexAttribArrayStride &nbsp;
GlVertexAttribArrayStride = 0x8624
//GlVertexAttribArrayType &nbsp;
GlVertexAttribArrayType = 0x8625
//GlVertexAttribArrayNormalized &nbsp;
GlVertexAttribArrayNormalized = 0x886a
//GlVertexAttribArrayPointer &nbsp;
GlVertexAttribArrayPointer = 0x8645
//GlVertexAttribArrayBufferBinding &nbsp;
GlVertexAttribArrayBufferBinding = 0x889f
//GlCullFace passed to <code>enable</code>/<code>disable</code> to turn on/off culling. can also be used with <code>getparameter</code> to find the current culling method.
GlCullFace = 0x0b44
//GlFront passed to <code>cullface</code> to specify that only front faces should be culled.
GlFront = 0x0404
//GlBack passed to <code>cullface</code> to specify that only back faces should be culled.
GlBack = 0x0405
//GlFrontAndBack passed to&nbsp;<code>cullface</code> to specify that front and back faces should be culled.
GlFrontAndBack = 0x0408
//GlBlend passed to <code>enable</code>/<code>disable</code> to turn on/off blending. can also be used with <code>getparameter</code> to find the current blending method.
GlBlend = 0x0be2
//GlDepthTest passed to <code>enable</code>/<code>disable</code> to turn on/off the depth test. can also be used with <code>getparameter</code> to query the depth test.
GlDepthTest = 0x0b71
//GlDither passed to <code>enable</code>/<code>disable</code> to turn on/off dithering. can also be used with <code>getparameter</code> to find the current dithering method.
GlDither = 0x0bd0
//GlPolygonOffsetFill passed to <code>enable</code>/<code>disable</code> to turn on/off the polygon offset. useful for rendering hidden-line images, decals, and or solids with highlighted edges. can also be used with <code>getparameter</code> to query the scissor test.
GlPolygonOffsetFill = 0x8037
//GlSampleAlphaToCoverage passed to <code>enable</code>/<code>disable</code> to turn on/off the alpha to coverage. used in multi-sampling alpha channels.
GlSampleAlphaToCoverage = 0x809e
//GlSampleCoverage passed to <code>enable</code>/<code>disable</code> to turn on/off the sample coverage. used in multi-sampling.
GlSampleCoverage = 0x80a0
//GlScissorTest passed to <code>enable</code>/<code>disable</code> to turn on/off the scissor test. can also be used with <code>getparameter</code> to query the scissor test.
GlScissorTest = 0x0c11
//GlStencilTest passed to <code>enable</code>/<code>disable</code> to turn on/off the stencil test. can also be used with <code>getparameter</code> to query the stencil test.
GlStencilTest = 0x0b90
//GlNoError returned from <code>geterror</code>.
GlNoError = 0
//GlInvalidEnum returned from <code>geterror</code>.
GlInvalidEnum = 0x0500
//GlInvalidValue returned from <code>geterror</code>.
GlInvalidValue = 0x0501
//GlInvalidOperation returned from <code>geterror</code>.
GlInvalidOperation = 0x0502
//GlOutOfMemory returned from <code>geterror</code>.
GlOutOfMemory = 0x0505
//GlContextLostWebgl returned from <code>geterror</code>.
GlContextLostWebgl = 0x9242
//GlCw passed to <code>frontface</code> to specify the front face of a polygon is drawn in the clockwise direction
GlCw = 0x0900
//GlCcw passed to <code>frontface</code> to specify the front face of a polygon is drawn in the counter clockwise direction
GlCcw = 0x0901
//GlDontCare there is no preference for this behavior.
GlDontCare = 0x1100
//GlFastest the most efficient behavior should be used.
GlFastest = 0x1101
//GlNicest the most correct or the highest quality option should be used.
GlNicest = 0x1102
//GlGenerateMipmapHint hint for the quality of filtering when generating mipmap images with <a href="/en-us/docs/web/api/webglrenderingcontext/generatemipmap" title="the webglrenderingcontext.generatemipmap() method of the webgl api generates a set of mipmaps for a webgltexture object."><code>webglrenderingcontext.generatemipmap()</code></a>.
GlGenerateMipmapHint = 0x8192
//GlByte &nbsp;
GlByte = 0x1400
//GlShort &nbsp;
GlShort = 0x1402
//GlUnsignedShort &nbsp;
GlUnsignedShort = 0x1403
//GlInt &nbsp;
GlInt = 0x1404
//GlUnsignedInt &nbsp;
GlUnsignedInt = 0x1405
//GlFloat &nbsp;
GlFloat = 0x1406
//GlDepthComponent &nbsp;
GlDepthComponent = 0x1902
//GlAlpha &nbsp;
GlAlpha = 0x1906
//GlRGB &nbsp;
GlRGB = 0x1907
//GlRGBA &nbsp;
GlRGBA = 0x1908
//GlLuminance &nbsp;
GlLuminance = 0x1909
//GlLuminanceAlpha &nbsp;
GlLuminanceAlpha = 0x190a
//GlUnsignedByte &nbsp;
GlUnsignedByte = 0x1401
//GlUnsignedShort4444 &nbsp;
GlUnsignedShort4444 = 0x8033
//GlUnsignedShort5551 &nbsp;
GlUnsignedShort5551 = 0x8034
//GlUnsignedShort565 &nbsp;
GlUnsignedShort565 = 0x8363
//GlFragmentShader passed to <code>createshader</code> to define a fragment shader.
GlFragmentShader = 0x8b30
//GlVertexShader passed to <code>createshader</code> to define a vertex shader
GlVertexShader = 0x8b31
//GlCompileStatus passed to <code>getshaderparamter</code> to get the status of the compilation. returns false if the shader was not compiled. you can then query <code>getshaderinfolog</code> to find the exact error
GlCompileStatus = 0x8b81
//GlDeleteStatus passed to <code>getshaderparamter</code> to determine if a shader was deleted via <code>deleteshader</code>. returns true if it was, false otherwise.
GlDeleteStatus = 0x8b80
//GlLinkStatus passed to <code>getprogramparameter</code> after calling <code>linkprogram</code> to determine if a program was linked correctly. returns false if there were errors. use <code>getprograminfolog</code> to find the exact error.
GlLinkStatus = 0x8b82
//GlValidateStatus passed to <code>getprogramparameter</code> after calling <code>validateprogram</code> to determine if it is valid. returns false if errors were found.
GlValidateStatus = 0x8b83
//GlAttachedShaders passed to <code>getprogramparameter</code> after calling <code>attachshader</code> to determine if the shader was attached correctly. returns false if errors occurred.
GlAttachedShaders = 0x8b85
//GlActiveAttributes passed to <code>getprogramparameter</code> to get the number of attributes active in a program.
GlActiveAttributes = 0x8b89
//GlActiveUniforms passed to <code>getprogramparamter</code> to get the number of uniforms active in a program.
GlActiveUniforms = 0x8b86
//GlMaxVertexAttribs the maximum number of entries possible in the vertex attribute list.
GlMaxVertexAttribs = 0x8869
//GlMaxVertexUniformVectors &nbsp;
GlMaxVertexUniformVectors = 0x8dfb
//GlMaxVaryingVectors &nbsp;
GlMaxVaryingVectors = 0x8dfc
//GlMaxCombinedTextureImageUnits &nbsp;
GlMaxCombinedTextureImageUnits = 0x8b4d
//GlMaxVertexTextureImageUnits &nbsp;
GlMaxVertexTextureImageUnits = 0x8b4c
//GlMaxTextureImageUnits implementation dependent number of maximum texture units. at least 8.
GlMaxTextureImageUnits = 0x8872
//GlMaxFragmentUniformVectors &nbsp;
GlMaxFragmentUniformVectors = 0x8dfd
//GlShaderType &nbsp;
GlShaderType = 0x8b4f
//GlShadingLanguageVersion &nbsp;
GlShadingLanguageVersion = 0x8b8c
//GlCurrentProgram &nbsp;
GlCurrentProgram = 0x8b8d
//GlNever passed to <code>depthfunction</code> or <code>stencilfunction</code> to specify depth or stencil tests will never pass. i.e. nothing will be drawn.
GlNever = 0x0200
//GlLess passed to <code>depthfunction</code> or <code>stencilfunction</code> to specify depth or stencil tests will pass if the new depth value is less than the stored value.
GlLess = 0x0201
//GlEqual passed to <code>depthfunction</code> or <code>stencilfunction</code> to specify depth or stencil tests will pass if the new depth value is equals to the stored value.
GlEqual = 0x0202
//GlLEqual passed to <code>depthfunction</code> or <code>stencilfunction</code> to specify depth or stencil tests will pass if the new depth value is less than or equal to the stored value.
GlLEqual = 0x0203
//GlGreater passed to <code>depthfunction</code> or <code>stencilfunction</code> to specify depth or stencil tests will pass if the new depth value is greater than the stored value.
GlGreater = 0x0204
//GlNotEqual passed to <code>depthfunction</code> or <code>stencilfunction</code> to specify depth or stencil tests will pass if the new depth value is not equal to the stored value.
GlNotEqual = 0x0205
//GlGEqual passed to <code>depthfunction</code> or <code>stencilfunction</code> to specify depth or stencil tests will pass if the new depth value is greater than or equal to the stored value.
GlGEqual = 0x0206
//GlAlways passed to <code>depthfunction</code> or <code>stencilfunction</code> to specify depth or stencil tests will always pass. i.e. pixels will be drawn in the order they are drawn.
GlAlways = 0x0207
//GlKeep &nbsp;
GlKeep = 0x1e00
//GlReplace &nbsp;
GlReplace = 0x1e01
//GlIncr &nbsp;
GlIncr = 0x1e02
//GlDecr &nbsp;
GlDecr = 0x1e03
//GlInvert &nbsp;
GlInvert = 0x150a
//GlIncrWrap &nbsp;
GlIncrWrap = 0x8507
//GlDecrWrap &nbsp;
GlDecrWrap = 0x8508
//GlNearest &nbsp;
GlNearest = 0x2600
//GlLinear &nbsp;
GlLinear = 0x2601
//GlNearestMipmapNearest &nbsp;
GlNearestMipmapNearest = 0x2700
//GlLinearMipmapNearest &nbsp;
GlLinearMipmapNearest = 0x2701
//GlNearestMipmapLinear &nbsp;
GlNearestMipmapLinear = 0x2702
//GlLinearMipmapLinear &nbsp;
GlLinearMipmapLinear = 0x2703
//GlTextureMagFilter &nbsp;
GlTextureMagFilter = 0x2800
//GlTextureMinFilter &nbsp;
GlTextureMinFilter = 0x2801
//GlTextureWrapS &nbsp;
GlTextureWrapS = 0x2802
//GlTextureWrapT &nbsp;
GlTextureWrapT = 0x2803
//GlTexture2D &nbsp;
GlTexture2D = 0x0de1
//GlTexture &nbsp;
GlTexture = 0x1702
//GlTextureCubeMap &nbsp;
GlTextureCubeMap = 0x8513
//GlTextureBindingCubeMap &nbsp;
GlTextureBindingCubeMap = 0x8514
//GlTextureCubeMapPositiveX &nbsp;
GlTextureCubeMapPositiveX = 0x8515
//GlTextureCubeMapNegativeX &nbsp;
GlTextureCubeMapNegativeX = 0x8516
//GlTextureCubeMapPositiveY &nbsp;
GlTextureCubeMapPositiveY = 0x8517
//GlTextureCubeMapNegativeY &nbsp;
GlTextureCubeMapNegativeY = 0x8518
//GlTextureCubeMapPositiveZ &nbsp;
GlTextureCubeMapPositiveZ = 0x8519
//GlTextureCubeMapNegativeZ &nbsp;
GlTextureCubeMapNegativeZ = 0x851a
//GlMaxCubeMapTextureSize &nbsp;
GlMaxCubeMapTextureSize = 0x851c
//GlActiveTexture the current active texture unit.
GlActiveTexture = 0x84e0
//GlRepeat &nbsp;
GlRepeat = 0x2901
//GlClampToEdge &nbsp;
GlClampToEdge = 0x812f
//GlMirroredRepeat &nbsp;
GlMirroredRepeat = 0x8370
//GlFloatVec2 &nbsp;
GlFloatVec2 = 0x8b50
//GlFloatVec3 &nbsp;
GlFloatVec3 = 0x8b51
//GlFloatVec4 &nbsp;
GlFloatVec4 = 0x8b52
//GlIntVec2 &nbsp;
GlIntVec2 = 0x8b53
//GlIntVec3 &nbsp;
GlIntVec3 = 0x8b54
//GlIntVec4 &nbsp;
GlIntVec4 = 0x8b55
//GlBool &nbsp;
GlBool = 0x8b56
//GlBoolVec2 &nbsp;
GlBoolVec2 = 0x8b57
//GlBoolVec3 &nbsp;
GlBoolVec3 = 0x8b58
//GlBoolVec4 &nbsp;
GlBoolVec4 = 0x8b59
//GlFloatMat2 &nbsp;
GlFloatMat2 = 0x8b5a
//GlFloatMat3 &nbsp;
GlFloatMat3 = 0x8b5b
//GlFloatMat4 &nbsp;
GlFloatMat4 = 0x8b5c
//GlSampler2d &nbsp;
GlSampler2d = 0x8b5e
//GlSamplerCube &nbsp;
GlSamplerCube = 0x8b60
//GlLowFloat &nbsp;
GlLowFloat = 0x8df0
//GlMediumFloat &nbsp;
GlMediumFloat = 0x8df1
//GlHighFloat &nbsp;
GlHighFloat = 0x8df2
//GlLowInt &nbsp;
GlLowInt = 0x8df3
//GlMediumInt &nbsp;
GlMediumInt = 0x8df4
//GlHighInt &nbsp;
GlHighInt = 0x8df5
//GlFramebuffer &nbsp;
GlFramebuffer = 0x8d40
//GlRenderbuffer &nbsp;
GlRenderbuffer = 0x8d41
//GlRGBA4 &nbsp;
GlRGBA4 = 0x8056
//GlRGB5A1 &nbsp;
GlRGB5A1 = 0x8057
//GlRGB565 &nbsp;
GlRGB565 = 0x8d62
//GlDepthComponent16 &nbsp;
GlDepthComponent16 = 0x81a5
//GlStencilIndex8 &nbsp;
GlStencilIndex8 = 0x8d48
//GlDepthStencil &nbsp;
GlDepthStencil = 0x84f9
//GlRenderbufferWidth &nbsp;
GlRenderbufferWidth = 0x8d42
//GlRenderbufferHeight &nbsp;
GlRenderbufferHeight = 0x8d43
//GlRenderbufferInternalFormat &nbsp;
GlRenderbufferInternalFormat = 0x8d44
//GlRenderbufferRedSize &nbsp;
GlRenderbufferRedSize = 0x8d50
//GlRenderbufferGreenSize &nbsp;
GlRenderbufferGreenSize = 0x8d51
//GlRenderbufferBlueSize &nbsp;
GlRenderbufferBlueSize = 0x8d52
//GlRenderbufferAlphaSize &nbsp;
GlRenderbufferAlphaSize = 0x8d53
//GlRenderbufferDepthSize &nbsp;
GlRenderbufferDepthSize = 0x8d54
//GlRenderbufferStencilSize &nbsp;
GlRenderbufferStencilSize = 0x8d55
//GlFramebufferAttachmentObjectType &nbsp;
GlFramebufferAttachmentObjectType = 0x8cd0
//GlFramebufferAttachmentObjectName &nbsp;
GlFramebufferAttachmentObjectName = 0x8cd1
//GlFramebufferAttachmentTextureLevel &nbsp;
GlFramebufferAttachmentTextureLevel = 0x8cd2
//GlFramebufferAttachmentTextureCubeMapFace &nbsp;
GlFramebufferAttachmentTextureCubeMapFace = 0x8cd3
//GlColorAttachment0 &nbsp;
GlColorAttachment0 = 0x8ce0
//GlDepthAttachment &nbsp;
GlDepthAttachment = 0x8d00
//GlStencilAttachment &nbsp;
GlStencilAttachment = 0x8d20
//GlNone &nbsp;
GlNone = 0
//GlFramebufferComplete &nbsp;
GlFramebufferComplete = 0x8cd5
//GlFramebufferIncompleteAttachment &nbsp;
GlFramebufferIncompleteAttachment = 0x8cd6
//GlFramebufferIncompleteMissingAttachment &nbsp;
GlFramebufferIncompleteMissingAttachment = 0x8cd7
//GlFramebufferIncompleteDimensions &nbsp;
GlFramebufferIncompleteDimensions = 0x8cd9
//GlFramebufferUnsupported &nbsp;
GlFramebufferUnsupported = 0x8cdd
//GlFramebufferBinding &nbsp;
GlFramebufferBinding = 0x8ca6
//GlRenderbufferBinding &nbsp;
GlRenderbufferBinding = 0x8ca7
//GlMaxRenderbufferSize &nbsp;
GlMaxRenderbufferSize = 0x84e8
//GlInvalidFramebufferOperation &nbsp;
GlInvalidFramebufferOperation = 0x0506
//GlUnpackFlipYWebgl &nbsp;
GlUnpackFlipYWebgl = 0x9240
//GlUnpackPremultiplyAlphaWebgl &nbsp;
GlUnpackPremultiplyAlphaWebgl = 0x9241
//GlUnpackColorspaceConversionWebgl &nbsp;
GlUnpackColorspaceConversionWebgl = 0x9243
//GlReadBuffer &nbsp;
GlReadBuffer = 0x0c02
//GlUnpackRowLength &nbsp;
GlUnpackRowLength = 0x0cf2
//GlUnpackSkipRows &nbsp;
GlUnpackSkipRows = 0x0cf3
//GlUnpackSkipPixels &nbsp;
GlUnpackSkipPixels = 0x0cf4
//GlPackRowLength &nbsp;
GlPackRowLength = 0x0d02
//GlPackSkipRows &nbsp;
GlPackSkipRows = 0x0d03
//GlPackSkipPixels &nbsp;
GlPackSkipPixels = 0x0d04
//GlTextureBinding3d &nbsp;
GlTextureBinding3d = 0x806a
//GlUnpackSkipImages &nbsp;
GlUnpackSkipImages = 0x806d
//GlUnpackImageHeight &nbsp;
GlUnpackImageHeight = 0x806e
//GlMax3dTextureSize &nbsp;
GlMax3dTextureSize = 0x8073
//GlMaxElementsVertices &nbsp;
GlMaxElementsVertices = 0x80e8
//GlMaxElementsIndices &nbsp;
GlMaxElementsIndices = 0x80e9
//GlMaxTextureLodBias &nbsp;
GlMaxTextureLodBias = 0x84fd
//GlMaxFragmentUniformComponents &nbsp;
GlMaxFragmentUniformComponents = 0x8b49
//GlMaxVertexUniformComponents &nbsp;
GlMaxVertexUniformComponents = 0x8b4a
//GlMaxArrayTextureLayers &nbsp;
GlMaxArrayTextureLayers = 0x88ff
//GlMinProgramTexelOffset &nbsp;
GlMinProgramTexelOffset = 0x8904
//GlMaxProgramTexelOffset &nbsp;
GlMaxProgramTexelOffset = 0x8905
//GlMaxVaryingComponents &nbsp;
GlMaxVaryingComponents = 0x8b4b
//GlFragmentShaderDerivativeHint &nbsp;
GlFragmentShaderDerivativeHint = 0x8b8b
//GlRasterizerDiscard &nbsp;
GlRasterizerDiscard = 0x8c89
//GlVertexArrayBinding &nbsp;
GlVertexArrayBinding = 0x85b5
//GlMaxVertexOutputComponents &nbsp;
GlMaxVertexOutputComponents = 0x9122
//GlMaxFragmentInputComponents &nbsp;
GlMaxFragmentInputComponents = 0x9125
//GlMaxServerWaitTimeout &nbsp;
GlMaxServerWaitTimeout = 0x9111
//GlMaxElementIndex &nbsp;
GlMaxElementIndex = 0x8d6b
//GlRed &nbsp;
GlRed = 0x1903
//GlRGB8 &nbsp;
GlRGB8 = 0x8051
//GlRGBA8 &nbsp;
GlRGBA8 = 0x8058
//GlRGB10A2 &nbsp;
GlRGB10A2 = 0x8059
//GlTexture3d &nbsp;
GlTexture3d = 0x806f
//GlTextureWrapR &nbsp;
GlTextureWrapR = 0x8072
//GlTextureMinLod &nbsp;
GlTextureMinLod = 0x813a
//GlTextureMaxLod &nbsp;
GlTextureMaxLod = 0x813b
//GlTextureBaseLevel &nbsp;
GlTextureBaseLevel = 0x813c
//GlTextureMaxLevel &nbsp;
GlTextureMaxLevel = 0x813d
//GlTextureCompareMode &nbsp;
GlTextureCompareMode = 0x884c
//GlTextureCompareFunc &nbsp;
GlTextureCompareFunc = 0x884d
//GlSrgb &nbsp;
GlSrgb = 0x8c40
//GlSrgb8 &nbsp;
GlSrgb8 = 0x8c41
//GlSrgb8Alpha8 &nbsp;
GlSrgb8Alpha8 = 0x8c43
//GlCompareRefToTexture &nbsp;
GlCompareRefToTexture = 0x884e
//GlRGBA32f &nbsp;
GlRGBA32f = 0x8814
//GlRGB32f &nbsp;
GlRGB32f = 0x8815
//GlRGBA16f &nbsp;
GlRGBA16f = 0x881a
//GlRGB16f &nbsp;
GlRGB16f = 0x881b
//GlTexture2DArray &nbsp;
GlTexture2DArray = 0x8c1a
//GlTextureBinding2dArray &nbsp;
GlTextureBinding2dArray = 0x8c1d
//GlR11fG11fB10f &nbsp;
GlR11fG11fB10f = 0x8c3a
//GlRGB9E5 &nbsp;
GlRGB9E5 = 0x8c3d
//GlRGBA32ui &nbsp;
GlRGBA32ui = 0x8d70
//GlRGB32ui &nbsp;
GlRGB32ui = 0x8d71
//GlRGBA16ui &nbsp;
GlRGBA16ui = 0x8d76
//GlRGB16ui &nbsp;
GlRGB16ui = 0x8d77
//GlRGBA8ui &nbsp;
GlRGBA8ui = 0x8d7c
//GlRGB8ui &nbsp;
GlRGB8ui = 0x8d7d
//GlRGBA32i &nbsp;
GlRGBA32i = 0x8d82
//GlRGB32i &nbsp;
GlRGB32i = 0x8d83
//GlRGBA16i &nbsp;
GlRGBA16i = 0x8d88
//GlRGB16i &nbsp;
GlRGB16i = 0x8d89
//GlRGBA8i &nbsp;
GlRGBA8i = 0x8d8e
//GlRGB8i &nbsp;
GlRGB8i = 0x8d8f
//GlRedInteger &nbsp;
GlRedInteger = 0x8d94
//GlRGBInteger &nbsp;
GlRGBInteger = 0x8d98
//GlRGBAInteger &nbsp;
GlRGBAInteger = 0x8d99
//GlR8 &nbsp;
GlR8 = 0x8229
//GlRg8 &nbsp;
GlRg8 = 0x822b
//GlRGB10A2ui &nbsp;
GlRGB10A2ui = 0x906f
//GlTextureImmutableFormat &nbsp;
GlTextureImmutableFormat = 0x912f
//GlTextureImmutableLevels &nbsp;
GlTextureImmutableLevels = 0x82df
//GlUnsignedInt2101010Rev &nbsp;
GlUnsignedInt2101010Rev = 0x8368
//GlUnsignedInt10f11f11fRev &nbsp;
GlUnsignedInt10f11f11fRev = 0x8c3b
//GlUnsignedInt5999Rev &nbsp;
GlUnsignedInt5999Rev = 0x8c3e
//GlFloat32UnsignedInt248Rev &nbsp;
GlFloat32UnsignedInt248Rev = 0x8dad
//GlHalfFloat &nbsp;
GlHalfFloat = 0x140b
//GlRg &nbsp;
GlRg = 0x8227
//GlRgInteger &nbsp;
GlRgInteger = 0x8228
//GlInt2101010Rev &nbsp;
GlInt2101010Rev = 0x8d9f
//GlCurrentQuery &nbsp;
GlCurrentQuery = 0x8865
//GlQueryResult &nbsp;
GlQueryResult = 0x8866
//GlQueryResultAvailable &nbsp;
GlQueryResultAvailable = 0x8867
//GlAnySamplesPassed &nbsp;
GlAnySamplesPassed = 0x8c2f
//GlAnySamplesPassedConservative &nbsp;
GlAnySamplesPassedConservative = 0x8d6a
//GlMaxDrawBuffers &nbsp;
GlMaxDrawBuffers = 0x8824
//GlDrawBuffer0 &nbsp;
GlDrawBuffer0 = 0x8825
//GlDrawBuffer1 &nbsp;
GlDrawBuffer1 = 0x8826
//GlDrawBuffer2 &nbsp;
GlDrawBuffer2 = 0x8827
//GlDrawBuffer3 &nbsp;
GlDrawBuffer3 = 0x8828
//GlDrawBuffer4 &nbsp;
GlDrawBuffer4 = 0x8829
//GlDrawBuffer5 &nbsp;
GlDrawBuffer5 = 0x882a
//GlDrawBuffer6 &nbsp;
GlDrawBuffer6 = 0x882b
//GlDrawBuffer7 &nbsp;
GlDrawBuffer7 = 0x882c
//GlDrawBuffer8 &nbsp;
GlDrawBuffer8 = 0x882d
//GlDrawBuffer9 &nbsp;
GlDrawBuffer9 = 0x882e
//GlDrawBuffer10 &nbsp;
GlDrawBuffer10 = 0x882f
//GlDrawBuffer11 &nbsp;
GlDrawBuffer11 = 0x8830
//GlDrawBuffer12 &nbsp;
GlDrawBuffer12 = 0x8831
//GlDrawBuffer13 &nbsp;
GlDrawBuffer13 = 0x8832
//GlDrawBuffer14 &nbsp;
GlDrawBuffer14 = 0x8833
//GlDrawBuffer15 &nbsp;
GlDrawBuffer15 = 0x8834
//GlMaxColorAttachments &nbsp;
GlMaxColorAttachments = 0x8cdf
//GlColorAttachment1 &nbsp;
GlColorAttachment1 = 0x8ce1
//GlColorAttachment2 &nbsp;
GlColorAttachment2 = 0x8ce2
//GlColorAttachment3 &nbsp;
GlColorAttachment3 = 0x8ce3
//GlColorAttachment4 &nbsp;
GlColorAttachment4 = 0x8ce4
//GlColorAttachment5 &nbsp;
GlColorAttachment5 = 0x8ce5
//GlColorAttachment6 &nbsp;
GlColorAttachment6 = 0x8ce6
//GlColorAttachment7 &nbsp;
GlColorAttachment7 = 0x8ce7
//GlColorAttachment8 &nbsp;
GlColorAttachment8 = 0x8ce8
//GlColorAttachment9 &nbsp;
GlColorAttachment9 = 0x8ce9
//GlColorAttachment10 &nbsp;
GlColorAttachment10 = 0x8cea
//GlColorAttachment11 &nbsp;
GlColorAttachment11 = 0x8ceb
//GlColorAttachment12 &nbsp;
GlColorAttachment12 = 0x8cec
//GlColorAttachment13 &nbsp;
GlColorAttachment13 = 0x8ced
//GlColorAttachment14 &nbsp;
GlColorAttachment14 = 0x8cee
//GlColorAttachment15 &nbsp;
GlColorAttachment15 = 0x8cef
//GlSampler3d &nbsp;
GlSampler3d = 0x8b5f
//GlSampler2dShadow &nbsp;
GlSampler2dShadow = 0x8b62
//GlSampler2dArray &nbsp;
GlSampler2dArray = 0x8dc1
//GlSampler2dArrayShadow &nbsp;
GlSampler2dArrayShadow = 0x8dc4
//GlSamplerCubeShadow &nbsp;
GlSamplerCubeShadow = 0x8dc5
//GlIntSampler2d &nbsp;
GlIntSampler2d = 0x8dca
//GlIntSampler3d &nbsp;
GlIntSampler3d = 0x8dcb
//GlIntSamplerCube &nbsp;
GlIntSamplerCube = 0x8dcc
//GlIntSampler2dArray &nbsp;
GlIntSampler2dArray = 0x8dcf
//GlUnsignedIntSampler2d &nbsp;
GlUnsignedIntSampler2d = 0x8dd2
//GlUnsignedIntSampler3d &nbsp;
GlUnsignedIntSampler3d = 0x8dd3
//GlUnsignedIntSamplerCube &nbsp;
GlUnsignedIntSamplerCube = 0x8dd4
//GlUnsignedIntSampler2dArray &nbsp;
GlUnsignedIntSampler2dArray = 0x8dd7
//GlMaxSamples &nbsp;
GlMaxSamples = 0x8d57
//GlSamplerBinding &nbsp;
GlSamplerBinding = 0x8919
//GlPixelPackBuffer &nbsp;
GlPixelPackBuffer = 0x88eb
//GlPixelUnpackBuffer &nbsp;
GlPixelUnpackBuffer = 0x88ec
//GlPixelPackBufferBinding &nbsp;
GlPixelPackBufferBinding = 0x88ed
//GlPixelUnpackBufferBinding &nbsp;
GlPixelUnpackBufferBinding = 0x88ef
//GlCopyReadBuffer &nbsp;
GlCopyReadBuffer = 0x8f36
//GlCopyWriteBuffer &nbsp;
GlCopyWriteBuffer = 0x8f37
//GlCopyReadBufferBinding &nbsp;
GlCopyReadBufferBinding = 0x8f36
//GlCopyWriteBufferBinding &nbsp;
GlCopyWriteBufferBinding = 0x8f37
//GlFloatMat2x3 &nbsp;
GlFloatMat2x3 = 0x8b65
//GlFloatMat2x4 &nbsp;
GlFloatMat2x4 = 0x8b66
//GlFloatMat3x2 &nbsp;
GlFloatMat3x2 = 0x8b67
//GlFloatMat3x4 &nbsp;
GlFloatMat3x4 = 0x8b68
//GlFloatMat4x2 &nbsp;
GlFloatMat4x2 = 0x8b69
//GlFloatMat4x3 &nbsp;
GlFloatMat4x3 = 0x8b6a
//GlUnsignedIntVec2 &nbsp;
GlUnsignedIntVec2 = 0x8dc6
//GlUnsignedIntVec3 &nbsp;
GlUnsignedIntVec3 = 0x8dc7
//GlUnsignedIntVec4 &nbsp;
GlUnsignedIntVec4 = 0x8dc8
//GlUnsignedNormalized &nbsp;
GlUnsignedNormalized = 0x8c17
//GlSignedNormalized &nbsp;
GlSignedNormalized = 0x8f9c
//GlVertexAttribArrayInteger &nbsp;
GlVertexAttribArrayInteger = 0x88fd
//GlVertexAttribArrayDivisor &nbsp;
GlVertexAttribArrayDivisor = 0x88fe
//GlTransformFeedbackBufferMode &nbsp;
GlTransformFeedbackBufferMode = 0x8c7f
//GlMaxTransformFeedbackSeparateComponents &nbsp;
GlMaxTransformFeedbackSeparateComponents = 0x8c80
//GlTransformFeedbackVaryings &nbsp;
GlTransformFeedbackVaryings = 0x8c83
//GlTransformFeedbackBufferStart &nbsp;
GlTransformFeedbackBufferStart = 0x8c84
//GlTransformFeedbackBufferSize &nbsp;
GlTransformFeedbackBufferSize = 0x8c85
//GlTransformFeedbackPrimitivesWritten &nbsp;
GlTransformFeedbackPrimitivesWritten = 0x8c88
//GlMaxTransformFeedbackInterleavedComponents &nbsp;
GlMaxTransformFeedbackInterleavedComponents = 0x8c8a
//GlMaxTransformFeedbackSeparateAttribs &nbsp;
GlMaxTransformFeedbackSeparateAttribs = 0x8c8b
//GlInterleavedAttribs &nbsp;
GlInterleavedAttribs = 0x8c8c
//GlSeparateAttribs &nbsp;
GlSeparateAttribs = 0x8c8d
//GlTransformFeedbackBuffer &nbsp;
GlTransformFeedbackBuffer = 0x8c8e
//GlTransformFeedbackBufferBinding &nbsp;
GlTransformFeedbackBufferBinding = 0x8c8f
//GlTransformFeedback &nbsp;
GlTransformFeedback = 0x8e22
//GlTransformFeedbackPaused &nbsp;
GlTransformFeedbackPaused = 0x8e23
//GlTransformFeedbackActive &nbsp;
GlTransformFeedbackActive = 0x8e24
//GlTransformFeedbackBinding &nbsp;
GlTransformFeedbackBinding = 0x8e25
//GlFramebufferAttachmentColorEncoding &nbsp;
GlFramebufferAttachmentColorEncoding = 0x8210
//GlFramebufferAttachmentComponentType &nbsp;
GlFramebufferAttachmentComponentType = 0x8211
//GlFramebufferAttachmentRedSize &nbsp;
GlFramebufferAttachmentRedSize = 0x8212
//GlFramebufferAttachmentGreenSize &nbsp;
GlFramebufferAttachmentGreenSize = 0x8213
//GlFramebufferAttachmentBlueSize &nbsp;
GlFramebufferAttachmentBlueSize = 0x8214
//GlFramebufferAttachmentAlphaSize &nbsp;
GlFramebufferAttachmentAlphaSize = 0x8215
//GlFramebufferAttachmentDepthSize &nbsp;
GlFramebufferAttachmentDepthSize = 0x8216
//GlFramebufferAttachmentStencilSize &nbsp;
GlFramebufferAttachmentStencilSize = 0x8217
//GlFramebufferDefault &nbsp;
GlFramebufferDefault = 0x8218
//GlDepthStencilAttachment &nbsp;
GlDepthStencilAttachment = 0x821a
//GlDepth24Stencil8 &nbsp;
GlDepth24Stencil8 = 0x88f0
//GlDrawFramebufferBinding &nbsp;
GlDrawFramebufferBinding = 0x8ca6
//GlReadFramebuffer &nbsp;
GlReadFramebuffer = 0x8ca8
//GlDrawFramebuffer &nbsp;
GlDrawFramebuffer = 0x8ca9
//GlReadFramebufferBinding &nbsp;
GlReadFramebufferBinding = 0x8caa
//GlRenderbufferSamples &nbsp;
GlRenderbufferSamples = 0x8cab
//GlFramebufferAttachmentTextureLayer &nbsp;
GlFramebufferAttachmentTextureLayer = 0x8cd4
//GlFramebufferIncompleteMultisample &nbsp;
GlFramebufferIncompleteMultisample = 0x8d56
//GlUniformBuffer &nbsp;
GlUniformBuffer = 0x8a11
//GlUniformBufferBinding &nbsp;
GlUniformBufferBinding = 0x8a28
//GlUniformBufferStart &nbsp;
GlUniformBufferStart = 0x8a29
//GlUniformBufferSize &nbsp;
GlUniformBufferSize = 0x8a2a
//GlMaxVertexUniformBlocks &nbsp;
GlMaxVertexUniformBlocks = 0x8a2b
//GlMaxFragmentUniformBlocks &nbsp;
GlMaxFragmentUniformBlocks = 0x8a2d
//GlMaxCombinedUniformBlocks &nbsp;
GlMaxCombinedUniformBlocks = 0x8a2e
//GlMaxUniformBufferBindings &nbsp;
GlMaxUniformBufferBindings = 0x8a2f
//GlMaxUniformBlockSize &nbsp;
GlMaxUniformBlockSize = 0x8a30
//GlMaxCombinedVertexUniformComponents &nbsp;
GlMaxCombinedVertexUniformComponents = 0x8a31
//GlMaxCombinedFragmentUniformComponents &nbsp;
GlMaxCombinedFragmentUniformComponents = 0x8a33
//GlUniformBufferOffsetAlignment &nbsp;
GlUniformBufferOffsetAlignment = 0x8a34
//GlActiveUniformBlocks &nbsp;
GlActiveUniformBlocks = 0x8a36
//GlUniformType &nbsp;
GlUniformType = 0x8a37
//GlUniformSize &nbsp;
GlUniformSize = 0x8a38
//GlUniformBlockIndex &nbsp;
GlUniformBlockIndex = 0x8a3a
//GlUniformOffset &nbsp;
GlUniformOffset = 0x8a3b
//GlUniformArrayStride &nbsp;
GlUniformArrayStride = 0x8a3c
//GlUniformMatrixStride &nbsp;
GlUniformMatrixStride = 0x8a3d
//GlUniformIsRowMajor &nbsp;
GlUniformIsRowMajor = 0x8a3e
//GlUniformBlockBinding &nbsp;
GlUniformBlockBinding = 0x8a3f
//GlUniformBlockDataSize &nbsp;
GlUniformBlockDataSize = 0x8a40
//GlUniformBlockActiveUniforms &nbsp;
GlUniformBlockActiveUniforms = 0x8a42
//GlUniformBlockActiveUniformIndices &nbsp;
GlUniformBlockActiveUniformIndices = 0x8a43
//GlUniformBlockReferencedByVertexShader &nbsp;
GlUniformBlockReferencedByVertexShader = 0x8a44
//GlUniformBlockReferencedByFragmentShader &nbsp;
GlUniformBlockReferencedByFragmentShader = 0x8a46
//GlObjectType &nbsp;
GlObjectType = 0x9112
//GlSyncCondition &nbsp;
GlSyncCondition = 0x9113
//GlSyncStatus &nbsp;
GlSyncStatus = 0x9114
//GlSyncFlags &nbsp;
GlSyncFlags = 0x9115
//GlSyncFence &nbsp;
GlSyncFence = 0x9116
//GlSyncGpuCommandsComplete &nbsp;
GlSyncGpuCommandsComplete = 0x9117
//GlUnsignaled &nbsp;
GlUnsignaled = 0x9118
//GlSignaled &nbsp;
GlSignaled = 0x9119
//GlAlreadySignaled &nbsp;
GlAlreadySignaled = 0x911a
//GlTimeoutExpired &nbsp;
GlTimeoutExpired = 0x911b
//GlConditionSatisfied &nbsp;
GlConditionSatisfied = 0x911c
//GlWaitFailed &nbsp;
GlWaitFailed = 0x911d
//GlSyncFlushCommandsBit &nbsp;
GlSyncFlushCommandsBit = 0x00000001
//GlColor &nbsp;
GlColor = 0x1800
//GlStencil &nbsp;
GlStencil = 0x1802
//GlMin &nbsp;
GlMin = 0x8007
//GlDepthComponent24 &nbsp;
GlDepthComponent24 = 0x81a6
//GlStreamRead &nbsp;
GlStreamRead = 0x88e1
//GlStreamCopy &nbsp;
GlStreamCopy = 0x88e2
//GlStaticRead &nbsp;
GlStaticRead = 0x88e5
//GlStaticCopy &nbsp;
GlStaticCopy = 0x88e6
//GlDynamicRead &nbsp;
GlDynamicRead = 0x88e9
//GlDynamicCopy &nbsp;
GlDynamicCopy = 0x88ea
//GlDepthComponent32f &nbsp;
GlDepthComponent32f = 0x8cac
//GlDepth32fStencil8 &nbsp;
GlDepth32fStencil8 = 0x8cad
//GlInvalidIndex &nbsp;
GlInvalidIndex = 0xffffffff
//GlTimeoutIgnored &nbsp;
GlTimeoutIgnored = -1
//GlMaxClientWaitTimeoutWebgl &nbsp;
GlMaxClientWaitTimeoutWebgl = 0x9247
//GlVertexAttribArrayDivisorAngle describes the frequency divisor used for instanced rendering.
GlVertexAttribArrayDivisorAngle = 0x88fe
//GlUnmaskedVendorWebgl passed to <code>getparameter</code> to get the vendor string of the graphics driver.
GlUnmaskedVendorWebgl = 0x9245
//GlUnmaskedRendererWebgl passed to <code>getparameter</code> to get the renderer string of the graphics driver.
GlUnmaskedRendererWebgl = 0x9246
//GlMaxTextureMaxAnisotropyExt returns the maximum available anisotropy.
GlMaxTextureMaxAnisotropyExt = 0x84ff
//GlTextureMaxAnisotropyExt passed to <code>texparameter</code> to set the desired maximum anisotropy for a texture.
GlTextureMaxAnisotropyExt = 0x84fe
//GlCompressedRgbS3tcDxt1Ext a dxt1-compressed image in an rgb image format.
GlCompressedRgbS3tcDxt1Ext = 0x83f0
//GlCompressedRgbaS3tcDxt1Ext a dxt1-compressed image in an rgb image format with a simple on/off alpha value.
GlCompressedRgbaS3tcDxt1Ext = 0x83f1
//GlCompressedRgbaS3tcDxt3Ext a dxt3-compressed image in an rgba image format. compared to a 32-bit rgba texture, it offers 4:1 compression.
GlCompressedRgbaS3tcDxt3Ext = 0x83f2
//GlCompressedRgbaS3tcDxt5Ext a dxt5-compressed image in an rgba image format. it also provides a 4:1 compression, but differs to the dxt3 compression in how the alpha compression is done.
GlCompressedRgbaS3tcDxt5Ext = 0x83f3
//GlCompressedR11Eac one-channel (red) unsigned format compression.
GlCompressedR11Eac = 0x9270
//GlCompressedSignedR11Eac one-channel (red) signed format compression.
GlCompressedSignedR11Eac = 0x9271
//GlCompressedRg11Eac two-channel (red and green) unsigned format compression.
GlCompressedRg11Eac = 0x9272
//GlCompressedSignedRg11Eac two-channel (red and green) signed format compression.
GlCompressedSignedRg11Eac = 0x9273
//GlCompressedRgb8Etc2 compresses rbg8 data with no alpha channel.
GlCompressedRgb8Etc2 = 0x9274
//GlCompressedRgba8Etc2Eac compresses rgba8 data. the rgb part is encoded the same as <code>rgbEtc2</code>, but the alpha part is encoded separately.
GlCompressedRgba8Etc2Eac = 0x9275
//GlCompressedSrgb8Etc2 compresses srbg8 data with no alpha channel.
GlCompressedSrgb8Etc2 = 0x9276
//GlCompressedSrgb8Alpha8Etc2Eac compresses srgba8 data. the srgb part is encoded the same as <code>srgbEtc2</code>, but the alpha part is encoded separately.
GlCompressedSrgb8Alpha8Etc2Eac = 0x9277
//GlCompressedRgb8PunchthroughAlpha1Etc2 similar to <code>rgb8Etc</code>, but with ability to punch through the alpha channel, which means to make it completely opaque or transparent.
GlCompressedRgb8PunchthroughAlpha1Etc2 = 0x9278
//GlCompressedSrgb8PunchthroughAlpha1Etc2 similar to <code>srgb8Etc</code>, but with ability to punch through the alpha channel, which means to make it completely opaque or transparent.
GlCompressedSrgb8PunchthroughAlpha1Etc2 = 0x9279
//GlCompressedRgbPvrtc4bppv1Img rgb compression in 4-bit mode. one block for each 4×4 pixels.
GlCompressedRgbPvrtc4bppv1Img = 0x8c00
//GlCompressedRgbaPvrtc4bppv1Img rgba compression in 4-bit mode. one block for each 4×4 pixels.
GlCompressedRgbaPvrtc4bppv1Img = 0x8c02
//GlCompressedRgbPvrtc2bppv1Img rgb compression in 2-bit mode. one block for each 8×4 pixels.
GlCompressedRgbPvrtc2bppv1Img = 0x8c01
//GlCompressedRgbaPvrtc2bppv1Img rgba compression in 2-bit mode. one block for each 8×4 pixe
GlCompressedRgbaPvrtc2bppv1Img = 0x8c03
//GlCompressedRgbEtc1Webgl compresses 24-bit rgb data with no alpha channel.
GlCompressedRgbEtc1Webgl = 0x8d64
//GlCompressedRgbAtcWebgl compresses rgb textures with no alpha channel.
GlCompressedRgbAtcWebgl = 0x8c92
//GlCompressedRgbaAtcExplicitAlphaWebgl compresses rgba textures using explicit alpha encoding (useful when alpha transitions are sharp).
GlCompressedRgbaAtcExplicitAlphaWebgl = 0x8c92
//GlCompressedRgbaAtcInterpolatedAlphaWebgl compresses rgba textures using interpolated alpha encoding (useful when alpha transitions are gradient).
GlCompressedRgbaAtcInterpolatedAlphaWebgl = 0x87ee
//GlUnsignedInt248Webgl unsigned integer type for 24-bit depth texture data.
GlUnsignedInt248Webgl = 0x84fa
//GlHalfFloatOes half floating-point type (16-bit).
GlHalfFloatOes = 0x8d61
//GlRGBA32fExt rgba 32-bit floating-point&nbsp;color-renderable format.
GlRGBA32fExt = 0x8814
//GlRGB32fExt rgb 32-bit floating-point&nbsp;color-renderable format.
GlRGB32fExt = 0x8815
//GlFramebufferAttachmentComponentTypeExt &nbsp;
GlFramebufferAttachmentComponentTypeExt = 0x8211
//GlUnsignedNormalizedExt &nbsp;
GlUnsignedNormalizedExt = 0x8c17
//GlMinExt produces the minimum color components of the source and destination colors.
GlMinExt = 0x8007
//GlMaxExt produces the maximum color components of the source and destination colors.
GlMaxExt = 0x8008
//GlSrgbExt unsized srgb format that leaves the precision up to the driver.
GlSrgbExt = 0x8c40
//GlSrgbAlphaExt unsized srgb format with unsized alpha component.
GlSrgbAlphaExt = 0x8c42
//GlSrgb8Alpha8Ext sized (8-bit) srgb and alpha formats.
GlSrgb8Alpha8Ext = 0x8c43
//GlFramebufferAttachmentColorEncodingExt returns the framebuffer color encoding.
GlFramebufferAttachmentColorEncodingExt = 0x8210
//GlFragmentShaderDerivativeHintOes indicates the accuracy of the derivative calculation for the glsl built-in functions: <code>dfdx</code>, <code>dfdy</code>, and <code>fwidth</code>.
GlFragmentShaderDerivativeHintOes = 0x8b8b
//GlColorAttachment0Webgl framebuffer color attachment point
GlColorAttachment0Webgl = 0x8ce0
//GlColorAttachment1Webgl framebuffer color attachment point
GlColorAttachment1Webgl = 0x8ce1
//GlColorAttachment2Webgl framebuffer color attachment point
GlColorAttachment2Webgl = 0x8ce2
//GlColorAttachment3Webgl framebuffer color attachment point
GlColorAttachment3Webgl = 0x8ce3
//GlColorAttachment4Webgl framebuffer color attachment point
GlColorAttachment4Webgl = 0x8ce4
//GlColorAttachment5Webgl framebuffer color attachment point
GlColorAttachment5Webgl = 0x8ce5
//GlColorAttachment6Webgl framebuffer color attachment point
GlColorAttachment6Webgl = 0x8ce6
//GlColorAttachment7Webgl framebuffer color attachment point
GlColorAttachment7Webgl = 0x8ce7
//GlColorAttachment8Webgl framebuffer color attachment point
GlColorAttachment8Webgl = 0x8ce8
//GlColorAttachment9Webgl framebuffer color attachment point
GlColorAttachment9Webgl = 0x8ce9
//GlColorAttachment10Webgl framebuffer color attachment point
GlColorAttachment10Webgl = 0x8cea
//GlColorAttachment11Webgl framebuffer color attachment point
GlColorAttachment11Webgl = 0x8ceb
//GlColorAttachment12Webgl framebuffer color attachment point
GlColorAttachment12Webgl = 0x8cec
//GlColorAttachment13Webgl framebuffer color attachment point
GlColorAttachment13Webgl = 0x8ced
//GlColorAttachment14Webgl framebuffer color attachment point
GlColorAttachment14Webgl = 0x8cee
//GlColorAttachment15Webgl framebuffer color attachment point
GlColorAttachment15Webgl = 0x8cef
//GlDrawBuffer0Webgl draw buffer
GlDrawBuffer0Webgl = 0x8825
//GlDrawBuffer1Webgl draw buffer
GlDrawBuffer1Webgl = 0x8826
//GlDrawBuffer2Webgl draw buffer
GlDrawBuffer2Webgl = 0x8827
//GlDrawBuffer3Webgl draw buffer
GlDrawBuffer3Webgl = 0x8828
//GlDrawBuffer4Webgl draw buffer
GlDrawBuffer4Webgl = 0x8829
//GlDrawBuffer5Webgl draw buffer
GlDrawBuffer5Webgl = 0x882a
//GlDrawBuffer6Webgl draw buffer
GlDrawBuffer6Webgl = 0x882b
//GlDrawBuffer7Webgl draw buffer
GlDrawBuffer7Webgl = 0x882c
//GlDrawBuffer8Webgl draw buffer
GlDrawBuffer8Webgl = 0x882d
//GlDrawBuffer9Webgl draw buffer
GlDrawBuffer9Webgl = 0x882e
//GlDrawBuffer10Webgl draw buffer
GlDrawBuffer10Webgl = 0x882f
//GlDrawBuffer11Webgl draw buffer
GlDrawBuffer11Webgl = 0x8830
//GlDrawBuffer12Webgl draw buffer
GlDrawBuffer12Webgl = 0x8831
//GlDrawBuffer13Webgl draw buffer
GlDrawBuffer13Webgl = 0x8832
//GlDrawBuffer14Webgl draw buffer
GlDrawBuffer14Webgl = 0x8833
//GlDrawBuffer15Webgl draw buffer
GlDrawBuffer15Webgl = 0x8834
//GlMaxColorAttachmentsWebgl maximum number of framebuffer color attachment points
GlMaxColorAttachmentsWebgl = 0x8cdf
//GlMaxDrawBuffersWebgl maximum number of draw buffers
GlMaxDrawBuffersWebgl = 0x8824
//GlVertexArrayBindingOes the bound vertex array object (vao).
GlVertexArrayBindingOes = 0x85b5
//GlQueryCounterBitsExt the number of bits used to hold the query result for the given target.
GlQueryCounterBitsExt = 0x8864
//GlCurrentQueryExt the currently active query.
GlCurrentQueryExt = 0x8865
//GlQueryResultExt the query result.
GlQueryResultExt = 0x8866
//GlQueryResultAvailableExt a boolean indicating whether or not a query result is available.
GlQueryResultAvailableExt = 0x8867
//GlTimeElapsedExt elapsed time (in nanoseconds).
GlTimeElapsedExt = 0x88bf
//GlTimestampExt the current time.
GlTimestampExt = 0x8e28
//GlGpuDisjointExt a boolean indicating whether or not the gpu performed any disjoint operation.
GlGpuDisjointExt = 0x8fbb
)

const (
	//GlTexture0 A texture unit. The first texture unit
	GlTexture0 = iota + 0x84c0
	//GlTexture1 A texture unit
	GlTexture1
	//GlTexture2 A texture unit
	GlTexture2
	//GlTexture3 A texture unit
	GlTexture3
	//GlTexture4 A texture unit
	GlTexture4
	//GlTexture5 A texture unit
	GlTexture5
	//GlTexture6 A texture unit
	GlTexture6
	//GlTexture7 A texture unit
	GlTexture7
	//GlTexture8 A texture unit
	GlTexture8
	//GlTexture9 A texture unit
	GlTexture9
	//GlTexture10 A texture unit
	GlTexture10
	//GlTexture11 A texture unit
	GlTexture11
	//GlTexture12 A texture unit
	GlTexture12
	//GlTexture13 A texture unit
	GlTexture13
	//GlTexture14 A texture unit
	GlTexture14
	//GlTexture15 A texture unit
	GlTexture15
	//GlTexture16 A texture unit
	GlTexture16
	//GlTexture17 A texture unit
	GlTexture17
	//GlTexture18 A texture unit
	GlTexture18
	//GlTexture19 A texture unit
	GlTexture19
	//GlTexture20 A texture unit
	GlTexture20
	//GlTexture21 A texture unit
	GlTexture21
	//GlTexture22 A texture unit
	GlTexture22
	//GlTexture23 A texture unit
	GlTexture23
	//GlTexture24 A texture unit
	GlTexture24
	//GlTexture25 A texture unit
	GlTexture25
	//GlTexture26 A texture unit
	GlTexture26
	//GlTexture27 A texture unit
	GlTexture27
	//GlTexture28 A texture unit
	GlTexture28
	//GlTexture29 A texture unit
	GlTexture29
	//GlTexture30 A texture unit
	GlTexture30
	//GlTexture31 A texture unit. The last texture unit.
	GlTexture31 = 0x84df
)
