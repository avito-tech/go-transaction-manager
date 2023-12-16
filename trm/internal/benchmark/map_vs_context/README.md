# Result

Storing transaction in context is faster than map storing. However, when using the real database (sqlite), the result is similar.

```
BenchmarkContextEmptyTransaction-12    	11344580	       100.8 ns/op
BenchmarkMapEmptyTransaction-12        	 1448832	       837.7 ns/op
BenchmarkContextCopy-12                	  622712	      1769 ns/op
BenchmarkMapCopy-12                    	  645990	      2562 ns/op
BenchmarkContextRealTransaction-12     	   12331	     98315 ns/op
BenchmarkMapRealTransaction-12         	   10000	    105360 ns/op
```