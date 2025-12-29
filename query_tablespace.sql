SELECT /* + RULE */
 df.tablespace_name "Tablespace",
 to_char(ROUND(nvl(df.MAXBYTES,0) / (1024 * 1024 * 1024), 2),9990.09) "MaxExtendGB",
 to_char(ROUND(nvl(df.bytes,0) / (1024 * 1024 * 1024), 2),9990.09) "FileSizeGB",
 to_char(Round(nvl((df.bytes - SUM(fs.bytes)),0) / (1024 * 1024 * 1024), 2),9990.09) "UsedGB",
 to_char(Round((nvl(df.bytes,0) - SUM(nvl(fs.bytes,0)))*100/nvl(df.bytes,0), 2),9990.09) "UsedFileSizePct",
 to_char(Round(nvl((df.bytes - SUM(fs.bytes)),0) / nvl(df.MAXBYTES,0) * 100 ,2),9990.09) "UsedMaxExtendPct"
  FROM dba_free_space fs,
       (SELECT i.tablespace_name,
        SUM(i.bytes) bytes,
        sum(case
              when i.autoextensible = 'YES' then
                case when i.MAXBYTES >= i.bytes then i.MAXBYTES else i.bytes end
              when i.autoextensible = 'NO' then
                case when i.bytes >= i.MAXBYTES  then i.bytes else i.MAXBYTES end
            end)  MAXBYTES
   FROM dba_data_files i
  GROUP BY i.tablespace_name) df
 WHERE fs.tablespace_name(+) = df.tablespace_name
 GROUP BY df.tablespace_name, df.bytes,df.MAXBYTES
 order by 6 desc;

