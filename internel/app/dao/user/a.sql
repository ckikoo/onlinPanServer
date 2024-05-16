SELECT 
    COALESCE(p.spaceSize, u.total_space) AS spaceSize
FROM 
    tb_user u
JOIN 
    tb_vip v ON u.user_id = v.user_id
JOIN 
    tb_package p ON v.vip_package_id = p.id
WHERE 
    CURDATE() BETWEEN v.active_from AND v.active_until
ORDER BY 
    COALESCE(p.spaceSize, u.total_space) DESC
LIMIT 1;